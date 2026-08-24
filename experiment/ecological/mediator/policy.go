package mediator

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Policy freezes the respondent-visible operations and budgets for one
// ecological task. Paths to source, evidence, caches, and controller artifacts
// are deliberately runtime inputs rather than policy fields.
type Policy struct {
	SchemaVersion       int            `json:"schema_version"`
	TaskID              string         `json:"task_id"`
	AllowedTools        []string       `json:"allowed_tools"`
	MaxCalls            int            `json:"max_calls"`
	MaxCallsByTool      map[string]int `json:"max_calls_by_tool"`
	MaxResultBytes      int            `json:"max_result_bytes"`
	MaxTotalResultBytes int            `json:"max_total_result_bytes"`
	MaxReadLines        int            `json:"max_read_lines"`
	MaxListDepth        int            `json:"max_list_depth"`
	MaxListEntries      int            `json:"max_list_entries"`
	MaxSearchMatches    int            `json:"max_search_matches"`
	MaxSearchPattern    int            `json:"max_search_pattern_bytes"`
	MaxTextFileBytes    int64          `json:"max_text_file_bytes"`
	TestTargets         []string       `json:"test_targets"`
}

func LoadPolicy(path string) (Policy, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, nil, err
	}
	var policy Policy
	if err := json.Unmarshal(data, &policy); err != nil {
		return Policy{}, nil, fmt.Errorf("parse tool policy: %w", err)
	}
	if err := policy.Validate(); err != nil {
		return Policy{}, nil, err
	}
	return policy, data, nil
}

func (p Policy) Validate() error {
	if p.SchemaVersion != 1 {
		return fmt.Errorf("unsupported tool policy schema_version %d", p.SchemaVersion)
	}
	if p.TaskID == "" {
		return fmt.Errorf("tool policy task_id is required")
	}
	if p.MaxCalls < 1 || p.MaxResultBytes < 1 || p.MaxTotalResultBytes < p.MaxResultBytes {
		return fmt.Errorf("invalid overall tool budgets")
	}
	if p.MaxReadLines < 1 || p.MaxListDepth < 1 || p.MaxListEntries < 1 || p.MaxSearchMatches < 1 || p.MaxSearchPattern < 1 || p.MaxTextFileBytes < 1 {
		return fmt.Errorf("invalid operation limit")
	}
	allowed := map[string]bool{}
	for _, tool := range p.AllowedTools {
		if tool == "" || allowed[tool] {
			return fmt.Errorf("allowed_tools contains an empty or duplicate value")
		}
		allowed[tool] = true
		if p.MaxCallsByTool[tool] < 1 {
			return fmt.Errorf("max_calls_by_tool[%q] must be positive", tool)
		}
	}
	for tool := range p.MaxCallsByTool {
		if !allowed[tool] {
			return fmt.Errorf("max_calls_by_tool contains unallowed tool %q", tool)
		}
	}
	for _, required := range []string{"source_list", "source_read", "source_search"} {
		if !allowed[required] {
			return fmt.Errorf("required read-only tool %q is not allowed", required)
		}
	}
	seenTargets := map[string]bool{}
	for _, target := range p.TestTargets {
		if target == "" || seenTargets[target] {
			return fmt.Errorf("test_targets contains an empty or duplicate value")
		}
		seenTargets[target] = true
	}
	if allowed["run_focused_test"] && len(p.TestTargets) == 0 {
		return fmt.Errorf("run_focused_test requires at least one test target")
	}
	return nil
}

// Digest returns a canonical digest independent of JSON object key order.
func (p Policy) Digest() (string, error) {
	clone := p
	clone.AllowedTools = append([]string(nil), p.AllowedTools...)
	clone.TestTargets = append([]string(nil), p.TestTargets...)
	sort.Strings(clone.AllowedTools)
	sort.Strings(clone.TestTargets)
	data, err := json.Marshal(clone)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}
