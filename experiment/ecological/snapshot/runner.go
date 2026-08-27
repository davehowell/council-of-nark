package snapshot

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type Runner struct {
	Root       string
	Ecological string
	Cache      string
	Work       string
	ConfigPath string
	ConfigData []byte
	Config     Config

	Attempt    string
	Controller string
	Source     string
	phase      string
	startedAt  time.Time
	provenance Provenance
}

type Provenance struct {
	SchemaVersion int              `json:"schema_version"`
	TaskID        string           `json:"task_id"`
	CreatedAt     string           `json:"created_at"`
	Controller    map[string]any   `json:"controller"`
	Upstream      map[string]any   `json:"upstream"`
	Archive       map[string]any   `json:"archive"`
	Licenses      []map[string]any `json:"licenses"`
	Sanitization  []map[string]any `json:"sanitization"`
	Source        map[string]any   `json:"source"`
	Closure       map[string]any   `json:"closure"`
	Validation    map[string]any   `json:"validation"`
	Isolation     map[string]any   `json:"isolation"`
}

type commandResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
}

func New(configPath string) (*Runner, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("ecological snapshots require macOS Seatbelt")
	}
	current, err := user.Current()
	if err != nil {
		return nil, err
	}
	if current.Uid == "0" {
		return nil, fmt.Errorf("refusing to prepare ecological snapshots as root")
	}
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, fmt.Errorf("find repository root: %w", err)
	}
	root := strings.TrimSpace(string(out))
	if !filepath.IsAbs(configPath) {
		configPath = filepath.Join(root, configPath)
	}
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}
	config, data, err := loadConfig(configPath)
	if err != nil {
		return nil, err
	}
	ecological := filepath.Join(root, "experiment", "ecological")
	return &Runner{
		Root: root, Ecological: ecological, Cache: filepath.Join(ecological, "cache"),
		Work: filepath.Join(ecological, "work"), ConfigPath: configPath, ConfigData: data, Config: config,
	}, nil
}

func (r *Runner) Run() (attempt string, err error) {
	r.startedAt = time.Now().UTC()
	for path, module := range map[string]string{r.Cache: "council.local/ecological-cache", r.Work: "council.local/ecological-work"} {
		if err = ensureRuntimeModuleBoundary(path, module); err != nil {
			return "", err
		}
	}
	r.Attempt, err = r.newAttemptPath()
	if err != nil {
		return "", err
	}
	r.Controller = filepath.Join(r.Attempt, "controller")
	r.Source = filepath.Join(r.Attempt, "source")
	if err = os.MkdirAll(r.Controller, 0o700); err != nil {
		return r.Attempt, err
	}
	if err = r.writeStatus("running", ""); err != nil {
		return r.Attempt, err
	}
	defer func() {
		if err != nil {
			_ = r.writeStatus("failed", err.Error())
		}
	}()

	steps := []struct {
		name string
		run  func() error
	}{
		{"record-controller", r.recordController},
		{"fetch-and-verify", r.fetchAndVerify},
		{"archive-parent", r.archiveParent},
		{"sanitize-source", r.sanitizeSource},
		{"prepare-validation", r.prepareValidation},
		{"freeze-closure", r.freezeClosure},
		{"offline-regression", r.runOfflineRegression},
		{"isolation-probe", r.runSourceIsolationProbe},
		{"finalize", r.finalize},
	}
	for _, step := range steps {
		r.phase = step.name
		if statusErr := r.writeStatus("running", ""); statusErr != nil {
			return r.Attempt, statusErr
		}
		if stepErr := step.run(); stepErr != nil {
			return r.Attempt, fmt.Errorf("%s: %w", step.name, stepErr)
		}
	}
	r.phase = "complete"
	if err = r.writeStatus("success", ""); err != nil {
		return r.Attempt, err
	}
	return r.Attempt, nil
}

func ensureRuntimeModuleBoundary(path, module string) error {
	if err := os.MkdirAll(path, 0o700); err != nil {
		return err
	}
	expected := []byte("module " + module + "\n\ngo 1.22\n")
	boundary := filepath.Join(path, "go.mod")
	if data, err := os.ReadFile(boundary); err == nil {
		if !bytes.Equal(data, expected) {
			return fmt.Errorf("unexpected runtime module boundary at %s", boundary)
		}
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	return atomicWrite(boundary, expected, 0o600)
}

func (r *Runner) newAttemptPath() (string, error) {
	if err := os.MkdirAll(r.Work, 0o700); err != nil {
		return "", err
	}
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", err
	}
	name := r.startedAt.Format("20060102T150405Z") + "-" + r.Config.TaskID + "-" + hex.EncodeToString(suffix[:])
	path := filepath.Join(r.Work, name)
	if err := os.Mkdir(path, 0o700); err != nil {
		return "", err
	}
	return path, nil
}

func (r *Runner) writeStatus(status, errorText string) error {
	value := map[string]any{
		"schema_version": 1, "task_id": r.Config.TaskID, "status": status,
		"phase": r.phase, "started_at": r.startedAt.Format(time.RFC3339Nano),
	}
	if status != "running" {
		value["finished_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	}
	if errorText != "" {
		value["error"] = errorText
	}
	return writeJSON(filepath.Join(r.Attempt, "status.json"), value, 0o600)
}

func (r *Runner) recordController() error {
	commit, err := capture(r.Root, nil, "git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	status, err := capture(r.Root, nil, "git", "status", "--porcelain=v1")
	if err != nil {
		return err
	}
	configRel, _ := filepath.Rel(r.Root, r.ConfigPath)
	r.provenance = Provenance{
		SchemaVersion: 1, TaskID: r.Config.TaskID, CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Controller: map[string]any{
			"commit": strings.TrimSpace(commit), "tree_dirty": strings.TrimSpace(status) != "",
			"config_path": filepath.ToSlash(configRel), "config_sha256": shaBytes(r.ConfigData),
		},
		Upstream: map[string]any{}, Archive: map[string]any{}, Source: map[string]any{},
		Closure: map[string]any{}, Validation: map[string]any{}, Isolation: map[string]any{},
	}
	return writeJSON(filepath.Join(r.Controller, "config.json"), r.Config, 0o600)
}

func capture(cwd string, env []string, name string, args ...string) (string, error) {
	command := exec.Command(name, args...)
	command.Dir = cwd
	if env != nil {
		command.Env = env
	}
	var stderr bytes.Buffer
	command.Stderr = &stderr
	out, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return string(out), nil
}

func runCommand(cwd string, env []string, name string, args ...string) commandResult {
	started := time.Now()
	command := exec.Command(name, args...)
	command.Dir = cwd
	if env != nil {
		command.Env = env
	}
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	exit := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exit = exitErr.ExitCode()
		} else {
			exit = -1
			stderr.WriteString("\nlauncher error: " + err.Error())
		}
	}
	return commandResult{ExitCode: exit, Stdout: stdout.String(), Stderr: stderr.String(), Duration: time.Since(started)}
}

func writeCommandLog(path string, command []string, result commandResult) error {
	value := map[string]any{
		"command": command, "exit_code": result.ExitCode,
		"duration_seconds": result.Duration.Seconds(), "stdout": result.Stdout, "stderr": result.Stderr,
	}
	return writeJSON(path, value, 0o600)
}

func writeJSON(path string, value any, mode os.FileMode) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return atomicWrite(path, data, mode)
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	name := file.Name()
	defer os.Remove(name)
	if _, err = file.Write(data); err == nil {
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if err := os.Chmod(name, mode); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func relativeTo(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}

func sortedEnvironment(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		out = append(out, key+"="+values[key])
	}
	return out
}
