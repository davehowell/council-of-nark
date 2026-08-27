package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ModuleRecord struct {
	Path     string              `json:"path"`
	Version  string              `json:"version"`
	Sum      string              `json:"sum"`
	Licenses []ModuleLicenseFile `json:"licenses"`
}

type ModuleLicenseFile struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type closurePaths struct {
	Toolchain       string
	ToolchainDigest string
	Modules         string
	ModulesDigest   string
	GoBinary        string
	ModuleRecords   []ModuleRecord
}

func (r *Runner) prepareValidation() error {
	validation := filepath.Join(r.Controller, "validation")
	parentRoot := filepath.Join(validation, "parent")
	evidenceRoot := filepath.Join(validation, "evidence")
	if err := extractTar(filepath.Join(r.Controller, "archives", "parent.tar"), parentRoot); err != nil {
		return err
	}
	evidenceArchive := filepath.Join(r.Controller, "archives", "evidence.tar")
	if err := r.writeArchive(r.Config.EvidenceCommit, evidenceArchive, "git-archive-evidence.json"); err != nil {
		return err
	}
	if err := extractTar(evidenceArchive, evidenceRoot); err != nil {
		return err
	}
	for _, path := range []string{"go.mod", "go.sum"} {
		parent, err := r.gitBlob(r.Config.ParentCommit, path)
		if err != nil {
			return err
		}
		evidence, err := r.gitBlob(r.Config.EvidenceCommit, path)
		if err != nil {
			return err
		}
		if !bytes.Equal(parent, evidence) {
			return fmt.Errorf("%s changed between parent and evidence; one closure cannot validate both", path)
		}
	}
	var applied []map[string]any
	for _, path := range r.Config.FocusedTest.EvidenceTestFiles {
		parentPath := filepath.Join(parentRoot, filepath.FromSlash(path))
		_, statErr := os.Lstat(parentPath)
		if r.Config.FocusedTest.RequireAbsentAtParent && !os.IsNotExist(statErr) {
			return fmt.Errorf("evidence test %s was expected to be absent at parent", path)
		}
		data, err := r.gitBlob(r.Config.EvidenceCommit, path)
		if err != nil {
			return err
		}
		evidencePath := filepath.Join(evidenceRoot, filepath.FromSlash(path))
		evidenceDigest, err := shaFile(evidencePath)
		if err != nil {
			return err
		}
		if evidenceDigest != shaBytes(data) {
			return fmt.Errorf("evidence archive blob mismatch for %s", path)
		}
		if err := copyBlobTo(parentPath, data, 0o644); err != nil {
			return err
		}
		applied = append(applied, map[string]any{"path": path, "sha256": evidenceDigest, "bytes": len(data)})
	}
	if err := normalizeSource(parentRoot, false); err != nil {
		return err
	}
	if err := normalizeSource(evidenceRoot, false); err != nil {
		return err
	}
	r.provenance.Validation = map[string]any{
		"evidence_tests_applied_to_parent": applied,
		"package":                          r.Config.FocusedTest.Package, "run": r.Config.FocusedTest.Run,
	}
	return nil
}

func (r *Runner) freezeClosure() error {
	paths, err := r.freezeGoToolchain()
	if err != nil {
		return err
	}
	modules, records, digest, err := r.prefetchModules(paths)
	if err != nil {
		return err
	}
	paths.Modules, paths.ModulesDigest, paths.ModuleRecords = modules, digest, records
	if err := writeJSON(filepath.Join(r.Controller, "closure.json"), map[string]any{
		"toolchain_path": relativeTo(r.Ecological, paths.Toolchain), "toolchain_tree_sha256": paths.ToolchainDigest,
		"module_path": relativeTo(r.Ecological, paths.Modules), "module_tree_sha256": paths.ModulesDigest,
		"modules": records,
	}, 0o600); err != nil {
		return err
	}
	r.provenance.Closure = map[string]any{
		"toolchain_tree_sha256": paths.ToolchainDigest, "module_tree_sha256": paths.ModulesDigest,
		"go_version": r.Config.Toolchain.RequiredVersion, "modules": records,
		"network_policy": "network allowed only during controller-side module prefetch; validation uses GOPROXY=off, GOSUMDB=off, and deny-by-default Seatbelt",
	}
	return nil
}

func (r *Runner) freezeGoToolchain() (closurePaths, error) {
	version, err := capture(r.Root, nil, "go", "version")
	if err != nil {
		return closurePaths{}, err
	}
	version = strings.TrimSpace(version)
	if version != r.Config.Toolchain.RequiredVersion {
		return closurePaths{}, fmt.Errorf("Go version is %q, expected %q", version, r.Config.Toolchain.RequiredVersion)
	}
	gorootText, err := capture(r.Root, append(os.Environ(), "GOENV=off"), "go", "env", "GOROOT")
	if err != nil {
		return closurePaths{}, err
	}
	goroot := strings.TrimSpace(gorootText)
	staging := filepath.Join(r.Cache, "toolchains", ".staging-"+filepath.Base(r.Attempt))
	if err := os.MkdirAll(staging, 0o700); err != nil {
		return closurePaths{}, err
	}
	rootCopy := filepath.Join(staging, "root")
	if err := copyTree(goroot, rootCopy, true); err != nil {
		return closurePaths{}, err
	}
	notices := filepath.Join(staging, "notices")
	if err := os.MkdirAll(notices, 0o755); err != nil {
		return closurePaths{}, err
	}
	licenseSource := filepath.Join(goroot, "LICENSE")
	if _, err := os.Stat(licenseSource); os.IsNotExist(err) {
		licenseSource = filepath.Join(filepath.Dir(goroot), "LICENSE")
	}
	licenseData, err := os.ReadFile(licenseSource)
	if err != nil {
		return closurePaths{}, fmt.Errorf("locate Go toolchain license: %w", err)
	}
	if !bytes.Contains(licenseData, []byte("Redistribution and use in source and binary forms")) {
		return closurePaths{}, fmt.Errorf("Go toolchain license marker is absent")
	}
	if err := atomicWrite(filepath.Join(notices, "GO-LICENSE"), licenseData, 0o444); err != nil {
		return closurePaths{}, err
	}
	if err := normalizeSource(staging, true); err != nil {
		return closurePaths{}, err
	}
	manifest, err := treeManifest(staging)
	if err != nil {
		return closurePaths{}, err
	}
	versionFields := strings.Fields(version)
	if len(versionFields) != 4 {
		return closurePaths{}, fmt.Errorf("unexpected go version output %q", version)
	}
	key := versionFields[2] + "-" + strings.ReplaceAll(versionFields[3], "/", "-") + "-" + manifest.TreeSHA256[:16]
	target := filepath.Join(r.Cache, "toolchains", key)
	if _, err := os.Stat(target); os.IsNotExist(err) {
		if err := os.Rename(staging, target); err != nil {
			return closurePaths{}, err
		}
		if err := writeClosureMetadata(target, manifest, map[string]any{
			"kind": "go-toolchain", "version": version, "license_spdx": r.Config.Toolchain.LicenseSPDX,
		}); err != nil {
			return closurePaths{}, err
		}
	} else if err != nil {
		return closurePaths{}, err
	} else {
		existing, err := treeManifestExcludingMetadata(target)
		if err != nil {
			return closurePaths{}, err
		}
		if existing.TreeSHA256 != manifest.TreeSHA256 {
			return closurePaths{}, fmt.Errorf("existing toolchain closure %s failed digest verification", target)
		}
		if err := makeTreeWritable(staging); err != nil {
			return closurePaths{}, err
		}
		if err := os.RemoveAll(staging); err != nil {
			return closurePaths{}, err
		}
	}
	goBinary := filepath.Join(target, "root", "bin", "go")
	if info, err := os.Stat(goBinary); err != nil || info.Mode()&0o111 == 0 {
		return closurePaths{}, fmt.Errorf("frozen Go executable is unavailable: %v", err)
	}
	return closurePaths{Toolchain: target, ToolchainDigest: manifest.TreeSHA256, GoBinary: goBinary}, nil
}

func writeClosureMetadata(root string, manifest TreeManifest, metadata map[string]any) error {
	if err := os.Chmod(root, 0o755); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(root, "manifest.json"), manifest, 0o444); err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(root, "metadata.json"), metadata, 0o444); err != nil {
		return err
	}
	return os.Chmod(root, 0o555)
}

// BuildClosureManifest hashes the immutable root/notices/gomod payload while
// excluding cache metadata files, matching the exporter closure digest.
func BuildClosureManifest(root string) (TreeManifest, error) {
	return treeManifestExcludingMetadata(root)
}

func treeManifestExcludingMetadata(root string) (TreeManifest, error) {
	var records []FileRecord
	var byteCount int64
	for _, name := range []string{"root", "notices", "gomod"} {
		subtree := filepath.Join(root, name)
		info, err := os.Stat(subtree)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return TreeManifest{}, err
		}
		manifest, err := treeManifest(subtree)
		if err != nil {
			return TreeManifest{}, err
		}
		records = append(records, FileRecord{Path: name, Type: "directory", Mode: uint32(info.Mode().Perm())})
		for _, record := range manifest.Files {
			record.Path = name + "/" + record.Path
			records = append(records, record)
		}
		byteCount += manifest.ByteCount
	}
	sort.Slice(records, func(i, j int) bool { return records[i].Path < records[j].Path })
	canonical, err := json.Marshal(records)
	if err != nil {
		return TreeManifest{}, err
	}
	return TreeManifest{SchemaVersion: 1, TreeSHA256: shaBytes(canonical), FileCount: len(records), ByteCount: byteCount, Files: records}, nil
}

func (r *Runner) prefetchModules(paths closurePaths) (string, []ModuleRecord, string, error) {
	stagingRoot := filepath.Join(r.Controller, "closure-staging")
	gomod := filepath.Join(stagingRoot, "gomod")
	gocache := filepath.Join(stagingRoot, "gocache")
	if err := os.MkdirAll(gomod, 0o755); err != nil {
		return "", nil, "", err
	}
	env := goEnvironment(paths, gomod, gocache, true)
	evidenceRoot := filepath.Join(r.Controller, "validation", "evidence")
	args := []string{"list", "-deps", "-test", "-json", r.Config.FocusedTest.Package}
	result := runCommand(evidenceRoot, env, paths.GoBinary, args...)
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "module-prefetch.json"), append([]string{"<frozen-go>"}, args...), result); err != nil {
		return "", nil, "", err
	}
	if result.ExitCode != 0 {
		return "", nil, "", fmt.Errorf("module prefetch failed: %s", strings.TrimSpace(result.Stderr))
	}
	if err := atomicWrite(filepath.Join(r.Controller, "dependency-packages-prefetch.jsonstream"), []byte(result.Stdout), 0o600); err != nil {
		return "", nil, "", err
	}
	records, err := parseModules([]byte(result.Stdout), gomod)
	if err != nil {
		return "", nil, "", err
	}
	verify := runCommand(evidenceRoot, env, paths.GoBinary, "mod", "verify")
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "go-mod-verify-prefetch.json"), []string{"<frozen-go>", "mod", "verify"}, verify); err != nil {
		return "", nil, "", err
	}
	if verify.ExitCode != 0 {
		return "", nil, "", fmt.Errorf("go mod verify failed after prefetch: %s", strings.TrimSpace(verify.Stderr))
	}
	if err := os.RemoveAll(gocache); err != nil {
		return "", nil, "", err
	}
	if err := normalizeSource(gomod, true); err != nil {
		return "", nil, "", err
	}
	manifest, err := treeManifestExcludingMetadata(stagingRoot)
	if err != nil {
		return "", nil, "", err
	}
	key := r.Config.TaskID + "-" + manifest.TreeSHA256[:16]
	target := filepath.Join(r.Cache, "modules", key)
	if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
		return "", nil, "", err
	}
	if _, err := os.Stat(target); os.IsNotExist(err) {
		cacheStaging := filepath.Join(filepath.Dir(target), ".staging-"+filepath.Base(r.Attempt))
		if err := os.Mkdir(cacheStaging, 0o700); err != nil {
			return "", nil, "", err
		}
		// macOS refuses to move a read-only directory into an existing directory.
		// Temporarily open only the module-cache root, then restore it before hashing.
		if err := os.Chmod(gomod, 0o755); err != nil {
			return "", nil, "", err
		}
		if err := os.Rename(gomod, filepath.Join(cacheStaging, "gomod")); err != nil {
			return "", nil, "", err
		}
		if err := os.Chmod(filepath.Join(cacheStaging, "gomod"), 0o555); err != nil {
			return "", nil, "", err
		}
		combined, err := treeManifestExcludingMetadata(cacheStaging)
		if err != nil {
			return "", nil, "", err
		}
		if combined.TreeSHA256 != manifest.TreeSHA256 {
			return "", nil, "", fmt.Errorf("module closure digest changed while freezing")
		}
		if err := writeClosureMetadata(cacheStaging, combined, map[string]any{"kind": "go-module-cache", "modules": records}); err != nil {
			return "", nil, "", err
		}
		if err := os.Rename(cacheStaging, target); err != nil {
			return "", nil, "", err
		}
	} else if err != nil {
		return "", nil, "", err
	} else {
		existing, err := treeManifestExcludingMetadata(target)
		if err != nil {
			return "", nil, "", err
		}
		if existing.TreeSHA256 != manifest.TreeSHA256 {
			return "", nil, "", fmt.Errorf("existing module closure failed digest verification")
		}
		if err := makeTreeWritable(gomod); err != nil {
			return "", nil, "", err
		}
		if err := os.RemoveAll(gomod); err != nil {
			return "", nil, "", err
		}
	}
	return target, records, manifest.TreeSHA256, nil
}

func goEnvironment(paths closurePaths, gomod, gocache string, network bool) []string {
	values := map[string]string{
		"HOME": filepath.Join(gocache, "home"), "TMPDIR": filepath.Join(gocache, "tmp"),
		"GOCACHE": gocache, "GOMODCACHE": gomod, "GOROOT": filepath.Join(paths.Toolchain, "root"),
		"GOENV": "off", "GOWORK": "off", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0",
		"GOFLAGS": "-buildvcs=false", "GOSUMDB": "sum.golang.org", "GOPROXY": "https://proxy.golang.org,direct",
		"PATH": filepath.Join(paths.Toolchain, "root", "bin") + ":/usr/bin:/bin", "LANG": "en_US.UTF-8",
	}
	if !network {
		values["GOSUMDB"] = "off"
		values["GOPROXY"] = "off"
	}
	_ = os.MkdirAll(values["TMPDIR"], 0o700)
	_ = os.MkdirAll(values["HOME"], 0o700)
	_ = network
	return sortedEnvironment(values)
}

type listedPackage struct {
	Standard bool `json:"Standard"`
	Module   *struct {
		Path    string `json:"Path"`
		Version string `json:"Version"`
		Dir     string `json:"Dir"`
		Sum     string `json:"Sum"`
		Main    bool   `json:"Main"`
	} `json:"Module"`
}

func parseModules(stream []byte, gomod string) ([]ModuleRecord, error) {
	decoder := json.NewDecoder(bytes.NewReader(stream))
	modules := map[string]ModuleRecord{}
	for {
		var pkg listedPackage
		if err := decoder.Decode(&pkg); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode go list stream: %w", err)
		}
		if pkg.Standard || pkg.Module == nil || pkg.Module.Main {
			continue
		}
		module := pkg.Module
		if module.Path == "" || module.Version == "" || module.Sum == "" || module.Dir == "" {
			return nil, fmt.Errorf("incomplete module metadata for %#v", module)
		}
		rel, err := filepath.Rel(gomod, module.Dir)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, fmt.Errorf("module %s resolved outside the fresh module cache", module.Path)
		}
		key := module.Path + "@" + module.Version
		if _, exists := modules[key]; exists {
			continue
		}
		licenses, err := moduleLicenses(module.Dir, filepath.ToSlash(rel))
		if err != nil {
			return nil, fmt.Errorf("module %s: %w", key, err)
		}
		modules[key] = ModuleRecord{Path: module.Path, Version: module.Version, Sum: module.Sum, Licenses: licenses}
	}
	out := make([]ModuleRecord, 0, len(modules))
	for _, module := range modules {
		out = append(out, module)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].Version < out[j].Version
		}
		return out[i].Path < out[j].Path
	})
	if len(out) == 0 {
		return nil, fmt.Errorf("focused package unexpectedly has no external module closure")
	}
	return out, nil
}

func moduleLicenses(moduleDir, moduleRelative string) ([]ModuleLicenseFile, error) {
	entries, err := os.ReadDir(moduleDir)
	if err != nil {
		return nil, err
	}
	var files []ModuleLicenseFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if !(strings.HasPrefix(name, "license") || strings.HasPrefix(name, "copying") || strings.HasPrefix(name, "notice") || strings.HasPrefix(name, "patents")) {
			continue
		}
		digest, err := shaFile(filepath.Join(moduleDir, entry.Name()))
		if err != nil {
			return nil, err
		}
		files = append(files, ModuleLicenseFile{Path: moduleRelative + "/" + entry.Name(), SHA256: digest})
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	if len(files) == 0 {
		return nil, fmt.Errorf("no root license/notice file found")
	}
	return files, nil
}
