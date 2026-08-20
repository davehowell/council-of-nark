package harness

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runtimeInventory(config Config) ([]map[string]any, error) {
	providers := config.Providers
	if len(providers) == 0 {
		providers = []Provider{config.Provider}
	}
	seen := map[string]bool{}
	rows := []map[string]any{}
	for _, provider := range providers {
		key := provider.Adapter + ":" + provider.Model
		if seen[key] || provider.Adapter == "mock" {
			continue
		}
		seen[key] = true
		command, err := commandFor(provider, "probe", "probe", map[string]any{"type": "object"})
		if err != nil {
			return nil, err
		}
		paths := append([]string{}, command.Executables...)
		if provider.Adapter == "pi" && len(command.Args) > 1 {
			paths = append(paths, command.Args[1])
		}
		artifacts := []map[string]any{}
		for _, path := range uniquePaths(paths) {
			digest, err := shaFile(path)
			if err != nil {
				return nil, err
			}
			artifacts = append(artifacts, map[string]any{"path": path, "sha256": digest})
		}
		version := "unavailable"
		if out, err := exec.Command(provider.Adapter, "--version").CombinedOutput(); err == nil {
			version = strings.TrimSpace(string(out))
		}
		rows = append(rows, map[string]any{"adapter": provider.Adapter, "model": provider.Model, "version": version, "artifacts": artifacts})
	}
	return rows, nil
}

func (h *Harness) FreezeRun(configArgument string) (string, error) {
	if err := currentUserCheck(); err != nil {
		return "", err
	}
	if err := h.sandboxProbe(false); err != nil {
		return "", err
	}
	if err := h.requireCleanTree(); err != nil {
		return "", err
	}
	audit := exec.Command("python3", "scripts/public_audit.py")
	audit.Dir = h.Root
	var stdout, stderr bytes.Buffer
	audit.Stdout = &stdout
	audit.Stderr = &stderr
	if err := audit.Run(); err != nil {
		return "", fmt.Errorf("public audit failed: %s%s", stdout.String(), stderr.String())
	}
	commit, err := h.sourceCommit()
	if err != nil {
		return "", err
	}
	configPath := configArgument
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(h.Root, configPath)
	}
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(h.Root, configPath)
	if err != nil || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("config must be inside the repository")
	}
	relative = filepath.ToSlash(relative)
	configBytes, err := h.gitBlob(commit, relative)
	if err != nil {
		return "", err
	}
	var config Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return "", err
	}
	out, err := h.git("ls-tree", "-r", "--name-only", commit, "--", "experiment", "justfile", "go.mod", "scripts/public_audit.py")
	if err != nil {
		return "", err
	}
	assets := []Asset{}
	for _, path := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if path == "" || strings.HasPrefix(path, "experiment/results/") {
			continue
		}
		data, err := h.gitBlob(commit, path)
		if err != nil {
			return "", err
		}
		assets = append(assets, Asset{Path: path, SHA256: shaBytes(data)})
	}
	runName := fmt.Sprintf("%s-%s-%s", slugTimestamp(), config.Label, commit[:10])
	run := filepath.Join(h.Experiment, "runs", runName)
	if _, err := os.Stat(run); err == nil {
		return "", fmt.Errorf("run already exists: %s", h.relative(run))
	}
	if err := os.MkdirAll(run, 0o755); err != nil {
		return "", err
	}
	runtimeRows, err := runtimeInventory(config)
	if err != nil {
		return "", err
	}
	manifest := Freeze{SchemaVersion: 2, CreatedAt: utcNow(), SourceCommit: commit, ConfigSource: relative, ConfigSHA256: shaBytes(configBytes), Assets: assets, ExternalRuntime: runtimeRows, Isolation: map[string]any{
		"git_worktree":      "one clean detached worktree per model call; prompt assembly only",
		"process":           "one fresh Seatbelt-confined process per model call",
		"seatbelt_required": true, "seatbelt_profile": "deny by default; empty cwd; ephemeral HOME; repository and worktree unavailable to child",
		"session_persistence": false, "tools": false, "project_context": false,
		"network":             "provider CLI transport permitted; local model tools disabled",
		"provider_side_state": "unobservable; server-side search/tools must be disabled separately",
	}, PublicAudit: map[string]any{"command": "python3 scripts/public_audit.py", "passed": true, "stdout": strings.TrimSpace(stdout.String()), "stderr": strings.TrimSpace(stderr.String())}}
	if err := writeJSON(filepath.Join(run, "freeze.json"), manifest); err != nil {
		return "", err
	}
	if err := atomicWrite(filepath.Join(run, "config.json"), configBytes, 0o644); err != nil {
		return "", err
	}
	return h.relative(run), nil
}
