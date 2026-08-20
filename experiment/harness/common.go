package harness

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

var terminalStatuses = map[string]bool{"success": true, "malformed": true, "error": true, "blocked": true}

type Harness struct {
	Root       string
	Experiment string
}

func New() (*Harness, error) {
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("unsupported operating system %q: the maintained harness requires macOS Seatbelt", runtime.GOOS)
	}
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, fmt.Errorf("find repository root: %w", err)
	}
	root := strings.TrimSpace(string(out))
	return &Harness{Root: root, Experiment: filepath.Join(root, "experiment")}, nil
}

func utcNow() string        { return time.Now().UTC().Format(time.RFC3339Nano) }
func slugTimestamp() string { return time.Now().UTC().Format("20060102T150405Z") }

func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(value); err != nil {
		return fmt.Errorf("decode %s: %w", path, err)
	}
	return nil
}

func atomicWrite(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer os.Remove(name)
	if _, err = f.Write(data); err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
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

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return atomicWrite(path, data, 0o644)
}

func shaBytes(value []byte) string {
	d := sha256.Sum256(value)
	return hex.EncodeToString(d[:])
}
func shaText(value string) string { return shaBytes([]byte(value)) }
func shaFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	d := sha256.New()
	if _, err := io.Copy(d, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(d.Sum(nil)), nil
}

func opaqueID(seed string, length int, parts ...any) string {
	values := []string{seed}
	for _, part := range parts {
		values = append(values, fmt.Sprint(part))
	}
	value := shaText(strings.Join(values, "\x1f"))
	if length > len(value) {
		length = len(value)
	}
	return value[:length]
}

func commandOutput(cwd string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = cwd
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w: %s", name, strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}
func (h *Harness) git(args ...string) ([]byte, error) { return commandOutput(h.Root, "git", args...) }
func (h *Harness) sourceCommit() (string, error) {
	out, err := h.git("rev-parse", "HEAD")
	return strings.TrimSpace(string(out)), err
}
func (h *Harness) requireCleanTree() error {
	out, err := h.git("status", "--porcelain=v1")
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(out)) != "" {
		return errors.New("refusing to freeze a dirty tree; commit or discard every source change first")
	}
	return nil
}
func (h *Harness) gitBlob(commit, path string) ([]byte, error) { return h.git("show", commit+":"+path) }

func (h *Harness) resolveRun(value string) (string, error) {
	path := value
	if !filepath.IsAbs(path) {
		path = filepath.Join(h.Root, path)
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	runs, _ := filepath.Abs(filepath.Join(h.Experiment, "runs"))
	rel, err := filepath.Rel(runs, path)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", fmt.Errorf("run path must be below %s", filepath.Join("experiment", "runs"))
	}
	if stat, err := os.Stat(path); err != nil || !stat.IsDir() {
		return "", fmt.Errorf("run does not exist: %s", path)
	}
	return path, nil
}
func (h *Harness) relative(path string) string {
	rel, err := filepath.Rel(h.Root, path)
	if err == nil && !strings.HasPrefix(rel, "..") {
		return filepath.ToSlash(rel)
	}
	return filepath.ToSlash(path)
}

func readCallRecord(run, callID string) (map[string]any, error) {
	path := filepath.Join(run, "calls", callID, "record.json")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	var value map[string]any
	if err := readJSON(path, &value); err != nil {
		return nil, err
	}
	return value, nil
}
func successfulResponse(run, callID string) (map[string]any, error) {
	record, err := readCallRecord(run, callID)
	if err != nil {
		return nil, err
	}
	if record == nil || stringValue(record["status"]) != "success" {
		return nil, fmt.Errorf("dependency %s has no successful response", callID)
	}
	return mapValue(record["parsed_response"]), nil
}

func stringValue(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
func intValue(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	case json.Number:
		n, _ := x.Int64()
		return int(n)
	case nil:
		return 0
	}
	var n int
	fmt.Sscan(fmt.Sprint(v), &n)
	return n
}
func floatValue(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case json.Number:
		n, _ := x.Float64()
		return n
	case nil:
		return 0
	}
	var n float64
	fmt.Sscan(fmt.Sprint(v), &n)
	return n
}
func boolValue(v any) bool          { b, _ := v.(bool); return b }
func mapValue(v any) map[string]any { m, _ := v.(map[string]any); return m }
func sliceValue(v any) []any        { a, _ := v.([]any); return a }
func stringSlice(v any) []string {
	a := sliceValue(v)
	out := make([]string, 0, len(a))
	for _, x := range a {
		out = append(out, stringValue(x))
	}
	return out
}
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func copyFile(source, destination string, mode os.FileMode) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return atomicWrite(destination, data, mode)
}
