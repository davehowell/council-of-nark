package respondent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	SchemaVersion            int      `json:"schema_version"`
	TaskID                   string   `json:"task_id"`
	SnapshotControllerCommit string   `json:"snapshot_controller_commit"`
	SourceTreeSHA256         string   `json:"source_tree_sha256"`
	ToolPolicy               string   `json:"tool_policy"`
	PiExtension              string   `json:"pi_extension"`
	Provider                 Provider `json:"provider"`
	ActiveTools              []string `json:"active_tools"`
	DoctorTimeoutSeconds     int      `json:"doctor_timeout_seconds"`
}

type Provider struct {
	Model    string `json:"model"`
	Thinking string `json:"thinking"`
}

func LoadConfig(root, path string) (Config, []byte, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, nil, err
	}
	var config Config
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&config); err != nil {
		return Config{}, nil, fmt.Errorf("parse respondent config: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return Config{}, nil, fmt.Errorf("respondent config contains trailing JSON")
	}
	if err := config.Validate(); err != nil {
		return Config{}, nil, err
	}
	return config, data, nil
}

func (c Config) Validate() error {
	if c.SchemaVersion != 1 || c.TaskID == "" || c.SnapshotControllerCommit == "" || len(c.SourceTreeSHA256) != 64 {
		return fmt.Errorf("respondent config identity is incomplete")
	}
	if c.ToolPolicy == "" || c.PiExtension == "" || c.Provider.Model == "" {
		return fmt.Errorf("respondent config paths/model are required")
	}
	if c.Provider.Thinking != "off" && c.Provider.Thinking != "minimal" && c.Provider.Thinking != "low" && c.Provider.Thinking != "medium" && c.Provider.Thinking != "high" && c.Provider.Thinking != "xhigh" && c.Provider.Thinking != "max" {
		return fmt.Errorf("invalid thinking level %q", c.Provider.Thinking)
	}
	if c.DoctorTimeoutSeconds < 5 {
		return fmt.Errorf("doctor timeout must be at least 5 seconds")
	}
	return validateActiveTools(c.ActiveTools)
}

func validateActiveTools(tools []string) error {
	expected := map[string]bool{
		"source_list": true, "source_read": true, "source_search": true,
		"run_focused_test": true, "submit_ecological_review": true,
	}
	if len(tools) != len(expected) {
		return fmt.Errorf("active_tools must contain exactly the five ecological tools")
	}
	for _, tool := range tools {
		if !expected[tool] {
			return fmt.Errorf("unexpected or duplicate active tool %q", tool)
		}
		delete(expected, tool)
	}
	return nil
}
