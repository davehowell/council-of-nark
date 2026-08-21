package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var commitPattern = regexp.MustCompile(`^[0-9a-f]{40}$`)
var taskPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type Config struct {
	SchemaVersion          int                `json:"schema_version"`
	TaskID                 string             `json:"task_id"`
	RepositoryURL          string             `json:"repository_url"`
	ParentCommit           string             `json:"parent_commit"`
	EvidenceCommit         string             `json:"evidence_commit"`
	ExpectedEvidenceParent string             `json:"expected_evidence_parent"`
	LicenseSPDX            string             `json:"license_spdx"`
	LicenseFiles           []ExpectedFile     `json:"license_files"`
	Sanitization           SanitizationConfig `json:"sanitization"`
	FocusedTest            FocusedTestConfig  `json:"focused_test"`
	Toolchain              ToolchainConfig    `json:"toolchain"`
}

type ExpectedFile struct {
	Path     string `json:"path"`
	SHA256   string `json:"sha256"`
	Contains string `json:"contains"`
}

type Removal struct {
	Path     string `json:"path"`
	Reason   string `json:"reason"`
	Required bool   `json:"required"`
}

type SanitizationConfig struct {
	Remove           []Removal `json:"remove"`
	ForbiddenContent []string  `json:"forbidden_content"`
}

type FocusedTestConfig struct {
	Package               string   `json:"package"`
	Run                   string   `json:"run"`
	EvidenceTestFiles     []string `json:"evidence_test_files"`
	RequireAbsentAtParent bool     `json:"require_absent_at_parent"`
	ParentFailureContains string   `json:"parent_failure_contains"`
}

type ToolchainConfig struct {
	Kind            string `json:"kind"`
	RequiredVersion string `json:"required_version"`
	LicenseSPDX     string `json:"license_spdx"`
}

func loadConfig(path string) (Config, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, nil, err
	}
	var config Config
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&config); err != nil {
		return Config{}, nil, fmt.Errorf("decode %s: %w", path, err)
	}
	if err := validateConfig(config); err != nil {
		return Config{}, nil, err
	}
	return config, data, nil
}

func validateConfig(c Config) error {
	if c.SchemaVersion != 1 {
		return fmt.Errorf("unsupported snapshot config schema_version %d", c.SchemaVersion)
	}
	if !taskPattern.MatchString(c.TaskID) {
		return fmt.Errorf("invalid task_id %q", c.TaskID)
	}
	if !strings.HasPrefix(c.RepositoryURL, "https://github.com/") || strings.Contains(c.RepositoryURL, "@") {
		return fmt.Errorf("repository_url must be an unauthenticated https://github.com URL")
	}
	for label, value := range map[string]string{
		"parent_commit": c.ParentCommit, "evidence_commit": c.EvidenceCommit,
		"expected_evidence_parent": c.ExpectedEvidenceParent,
	} {
		if !commitPattern.MatchString(value) {
			return fmt.Errorf("%s must be a full lowercase SHA-1 commit ID", label)
		}
	}
	if c.ParentCommit != c.ExpectedEvidenceParent {
		return fmt.Errorf("expected_evidence_parent must equal parent_commit for this pipeline")
	}
	if c.LicenseSPDX == "" || len(c.LicenseFiles) == 0 {
		return fmt.Errorf("license_spdx and license_files are required")
	}
	seen := map[string]bool{}
	for _, expected := range c.LicenseFiles {
		if err := validateRelativePath(expected.Path); err != nil {
			return fmt.Errorf("license file: %w", err)
		}
		if !regexp.MustCompile(`^[0-9a-f]{64}$`).MatchString(expected.SHA256) {
			return fmt.Errorf("license file %s has invalid sha256", expected.Path)
		}
		if expected.Contains == "" {
			return fmt.Errorf("license file %s requires a content marker", expected.Path)
		}
		if seen[expected.Path] {
			return fmt.Errorf("duplicate configured path %s", expected.Path)
		}
		seen[expected.Path] = true
	}
	for _, removal := range c.Sanitization.Remove {
		if err := validateRelativePath(removal.Path); err != nil {
			return fmt.Errorf("removal: %w", err)
		}
		if removal.Reason == "" {
			return fmt.Errorf("removal %s requires a reason", removal.Path)
		}
		if seen[removal.Path] {
			return fmt.Errorf("duplicate or overlapping configured path %s", removal.Path)
		}
		seen[removal.Path] = true
	}
	removalPaths := make([]string, 0, len(c.Sanitization.Remove))
	for _, removal := range c.Sanitization.Remove {
		removalPaths = append(removalPaths, removal.Path)
	}
	sort.Strings(removalPaths)
	for i, parent := range removalPaths {
		for _, child := range removalPaths[i+1:] {
			if strings.HasPrefix(child, parent+"/") {
				return fmt.Errorf("overlapping removals %s and %s", parent, child)
			}
		}
	}
	for _, needle := range c.Sanitization.ForbiddenContent {
		if len(needle) < 7 {
			return fmt.Errorf("forbidden_content values must be at least 7 bytes")
		}
	}
	if c.FocusedTest.Package == "" || strings.HasPrefix(c.FocusedTest.Package, "/") || strings.Contains(c.FocusedTest.Package, "..") {
		return fmt.Errorf("focused_test.package must be a relative Go package")
	}
	if c.FocusedTest.Run == "" {
		return fmt.Errorf("focused_test.run is required")
	}
	if _, err := regexp.Compile(c.FocusedTest.Run); err != nil {
		return fmt.Errorf("focused_test.run: %w", err)
	}
	if len(c.FocusedTest.EvidenceTestFiles) == 0 || c.FocusedTest.ParentFailureContains == "" {
		return fmt.Errorf("focused test evidence files and expected parent failure are required")
	}
	packageDir := strings.TrimPrefix(filepath.ToSlash(c.FocusedTest.Package), "./")
	for _, path := range c.FocusedTest.EvidenceTestFiles {
		if err := validateRelativePath(path); err != nil {
			return fmt.Errorf("evidence test file: %w", err)
		}
		if !strings.HasPrefix(path, packageDir+"/") || !strings.HasSuffix(path, "_test.go") {
			return fmt.Errorf("evidence test file %s must be a _test.go file in %s", path, packageDir)
		}
	}
	if c.Toolchain.Kind != "go" || c.Toolchain.RequiredVersion == "" || c.Toolchain.LicenseSPDX == "" {
		return fmt.Errorf("a pinned Go toolchain and license are required")
	}
	if runtime.GOOS != "darwin" {
		return fmt.Errorf("unsupported operating system %q: ecological snapshots require macOS Seatbelt", runtime.GOOS)
	}
	return nil
}

func validateRelativePath(value string) error {
	if value == "" || filepath.IsAbs(value) || filepath.ToSlash(filepath.Clean(value)) != value || value == "." || strings.HasPrefix(value, "../") {
		return fmt.Errorf("unsafe relative path %q", value)
	}
	return nil
}
