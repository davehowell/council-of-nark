package respondent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// RunConfig freezes one ecological respondent stage. Multi-stage arms are
// assembled by a separate controller; this config governs one isolated session.
type RunConfig struct {
	SchemaVersion            int         `json:"schema_version"`
	TaskID                   string      `json:"task_id"`
	Adapter                  string      `json:"adapter"`
	SnapshotControllerCommit string      `json:"snapshot_controller_commit"`
	SourceTreeSHA256         string      `json:"source_tree_sha256"`
	ToolPolicy               string      `json:"tool_policy"`
	PiExtension              string      `json:"pi_extension"`
	SystemPrompt             string      `json:"system_prompt"`
	Brief                    string      `json:"brief"`
	Provider                 Provider    `json:"provider"`
	ActiveTools              []string    `json:"active_tools"`
	Budgets                  RunBudgets  `json:"budgets"`
	Mock                     *MockConfig `json:"mock,omitempty"`
}

type RunBudgets struct {
	TimeoutSeconds   int `json:"timeout_seconds"`
	MaxProviderTurns int `json:"max_provider_turns"`
	MaxTotalTokens   int `json:"max_total_tokens"`
	MaxFinalBytes    int `json:"max_final_bytes"`
}

type MockConfig struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func LoadRunConfig(root, path string) (RunConfig, []byte, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return RunConfig{}, nil, err
	}
	var config RunConfig
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return RunConfig{}, nil, fmt.Errorf("parse ecological run config: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return RunConfig{}, nil, fmt.Errorf("ecological run config contains trailing JSON")
	}
	if err := config.Validate(); err != nil {
		return RunConfig{}, nil, err
	}
	return config, data, nil
}

func (c RunConfig) Validate() error {
	if c.SchemaVersion != 1 || c.TaskID == "" || c.SnapshotControllerCommit == "" || len(c.SourceTreeSHA256) != 64 {
		return fmt.Errorf("ecological run identity is incomplete")
	}
	if c.Adapter != "mock" && c.Adapter != "pi" {
		return fmt.Errorf("adapter must be mock or pi")
	}
	if c.ToolPolicy == "" || c.PiExtension == "" || c.SystemPrompt == "" || c.Brief == "" {
		return fmt.Errorf("tool, extension, and prompt paths are required")
	}
	if c.Provider.Model == "" {
		return fmt.Errorf("provider model is required")
	}
	if c.Provider.Thinking != "off" && c.Provider.Thinking != "minimal" && c.Provider.Thinking != "low" && c.Provider.Thinking != "medium" && c.Provider.Thinking != "high" && c.Provider.Thinking != "xhigh" && c.Provider.Thinking != "max" {
		return fmt.Errorf("invalid thinking level %q", c.Provider.Thinking)
	}
	if err := validateActiveTools(c.ActiveTools); err != nil {
		return err
	}
	if c.Budgets.TimeoutSeconds < 5 || c.Budgets.MaxProviderTurns < 1 || c.Budgets.MaxTotalTokens < 1 || c.Budgets.MaxFinalBytes < 256 {
		return fmt.Errorf("ecological run budgets are invalid")
	}
	if c.Adapter == "mock" {
		if c.Mock == nil || c.Mock.InputTokens < 0 || c.Mock.OutputTokens < 0 || c.Mock.InputTokens+c.Mock.OutputTokens < 1 {
			return fmt.Errorf("mock adapter requires positive deterministic usage")
		}
		if c.Mock.InputTokens+c.Mock.OutputTokens > c.Budgets.MaxTotalTokens {
			return fmt.Errorf("mock usage exceeds the token budget")
		}
	} else if c.Mock != nil {
		return fmt.Errorf("mock settings are not allowed for the pi adapter")
	}
	return nil
}
