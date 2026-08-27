package mediator

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type Request struct {
	RequestID string          `json:"request_id"`
	Tool      string          `json:"tool"`
	Arguments json.RawMessage `json:"arguments"`
}

type Response struct {
	RequestID string         `json:"request_id"`
	Sequence  int            `json:"sequence"`
	OK        bool           `json:"ok"`
	Output    string         `json:"output,omitempty"`
	Error     string         `json:"error,omitempty"`
	Truncated bool           `json:"truncated"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

type TranscriptEntry struct {
	SchemaVersion   int      `json:"schema_version"`
	Sequence        int      `json:"sequence"`
	StartedAt       string   `json:"started_at"`
	DurationSeconds float64  `json:"duration_seconds"`
	Request         Request  `json:"request"`
	Response        Response `json:"response"`
}

type TestRunner interface {
	Run(target string) (output string, metadata map[string]any, err error)
}

type TestRunnerFunc func(target string) (string, map[string]any, error)

func (f TestRunnerFunc) Run(target string) (string, map[string]any, error) { return f(target) }

type Session struct {
	Source     string
	Policy     Policy
	Tests      TestRunner
	Transcript io.Writer

	mu               sync.Mutex
	sequence         int
	callsByTool      map[string]int
	totalResultBytes int
	allowed          map[string]bool
	testTargets      map[string]bool
	transcriptErr    error
}

func NewSession(source string, policy Policy, tests TestRunner, transcript io.Writer) (*Session, error) {
	if err := policy.Validate(); err != nil {
		return nil, err
	}
	canonical, err := filepath.EvalSymlinks(source)
	if err != nil {
		return nil, fmt.Errorf("resolve source root: %w", err)
	}
	canonical, err = filepath.Abs(canonical)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(canonical)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("source root is not a directory")
	}
	if _, err := os.Stat(filepath.Join(canonical, ".git")); err == nil || !os.IsNotExist(err) {
		return nil, fmt.Errorf("source root must not contain .git")
	}
	allowed := map[string]bool{}
	for _, tool := range policy.AllowedTools {
		allowed[tool] = true
	}
	targets := map[string]bool{}
	for _, target := range policy.TestTargets {
		targets[target] = true
	}
	return &Session{
		Source: canonical, Policy: policy, Tests: tests, Transcript: transcript,
		callsByTool: map[string]int{}, allowed: allowed, testTargets: targets,
	}, nil
}

// Handle serializes requests so budget accounting and the transcript have one
// unambiguous order even when a provider emits parallel tool calls.
func (s *Session) Handle(request Request) Response {
	s.mu.Lock()
	defer s.mu.Unlock()
	started := time.Now().UTC()
	s.sequence++
	response := Response{RequestID: request.RequestID, Sequence: s.sequence, Truncated: false}

	if s.transcriptErr != nil {
		response.Error = "mediator transcript is unavailable; session is fail-closed"
	} else if request.RequestID == "" {
		response.Error = "request_id is required"
	} else if !s.allowed[request.Tool] {
		response.Error = fmt.Sprintf("tool %q is not allowed", request.Tool)
	} else if s.sequence > s.Policy.MaxCalls {
		response.Error = "overall tool-call budget exhausted"
	} else if s.callsByTool[request.Tool] >= s.Policy.MaxCallsByTool[request.Tool] {
		response.Error = fmt.Sprintf("tool-call budget exhausted for %s", request.Tool)
	} else {
		s.callsByTool[request.Tool]++
		output, metadata, err := s.execute(request)
		if err != nil {
			response.Error = err.Error()
		} else {
			remaining := s.Policy.MaxTotalResultBytes - s.totalResultBytes
			limit := s.Policy.MaxResultBytes
			if remaining < limit {
				limit = remaining
			}
			if limit <= 0 {
				response.Error = "total tool-result byte budget exhausted"
			} else {
				response.Output, response.Truncated = truncateUTF8(output, limit)
				s.totalResultBytes += len([]byte(response.Output))
				response.OK = true
				if metadata == nil {
					metadata = map[string]any{}
				}
				response.Metadata = metadata
				response.Metadata["result_bytes"] = len([]byte(response.Output))
				response.Metadata["total_result_bytes"] = s.totalResultBytes
				response.Metadata["calls_used"] = s.sequence
			}
		}
	}
	if s.transcriptErr == nil {
		if err := s.writeTranscript(TranscriptEntry{
			SchemaVersion: 1, Sequence: s.sequence, StartedAt: started.Format(time.RFC3339Nano),
			DurationSeconds: time.Since(started).Seconds(), Request: request, Response: response,
		}); err != nil {
			s.transcriptErr = err
			response.OK = false
			response.Output = ""
			response.Metadata = nil
			response.Error = "mediator transcript write failed; session is fail-closed"
		}
	}
	return response
}

func (s *Session) execute(request Request) (string, map[string]any, error) {
	switch request.Tool {
	case "source_list":
		var args listArguments
		if err := strictArguments(request.Arguments, &args); err != nil {
			return "", nil, err
		}
		return s.list(args)
	case "source_read":
		var args readArguments
		if err := strictArguments(request.Arguments, &args); err != nil {
			return "", nil, err
		}
		return s.read(args)
	case "source_search":
		var args searchArguments
		if err := strictArguments(request.Arguments, &args); err != nil {
			return "", nil, err
		}
		return s.search(args)
	case "run_focused_test":
		var args testArguments
		if err := strictArguments(request.Arguments, &args); err != nil {
			return "", nil, err
		}
		if !s.testTargets[args.Target] {
			return "", nil, fmt.Errorf("test target %q is not allowed", args.Target)
		}
		if s.Tests == nil {
			return "", nil, fmt.Errorf("focused test runner is unavailable")
		}
		return s.Tests.Run(args.Target)
	default:
		return "", nil, fmt.Errorf("tool %q has no implementation", request.Tool)
	}
}

func strictArguments(raw json.RawMessage, target any) error {
	if len(raw) == 0 {
		raw = []byte("{}")
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid tool arguments: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return fmt.Errorf("invalid tool arguments: trailing JSON value")
		}
		return fmt.Errorf("invalid tool arguments: %w", err)
	}
	return nil
}

func (s *Session) writeTranscript(entry TranscriptEntry) error {
	if s.Transcript == nil {
		return nil
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	frame := append(data, '\n')
	written, err := s.Transcript.Write(frame)
	if err == nil && written != len(frame) {
		return io.ErrShortWrite
	}
	return err
}

func (s *Session) resolve(relative string, requireDirectory bool) (string, string, error) {
	if relative == "" {
		relative = "."
	}
	if strings.ContainsRune(relative, '\x00') || filepath.IsAbs(relative) || filepath.VolumeName(relative) != "" {
		return "", "", fmt.Errorf("path must be relative to the source root")
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", "", fmt.Errorf("path escapes the source root")
	}
	if clean == ".git" || strings.HasPrefix(clean, ".git"+string(filepath.Separator)) {
		return "", "", fmt.Errorf("Git metadata is not available")
	}
	candidate := filepath.Join(s.Source, clean)
	current := s.Source
	if clean != "." {
		for _, component := range strings.Split(clean, string(filepath.Separator)) {
			current = filepath.Join(current, component)
			info, err := os.Lstat(current)
			if err != nil {
				return "", "", fmt.Errorf("path is unavailable: %s", filepath.ToSlash(clean))
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return "", "", fmt.Errorf("symbolic links are not available")
			}
		}
	}
	canonical, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		return "", "", fmt.Errorf("path is unavailable: %s", filepath.ToSlash(clean))
	}
	if canonical != s.Source && !strings.HasPrefix(canonical, s.Source+string(filepath.Separator)) {
		return "", "", fmt.Errorf("resolved path escapes the source root")
	}
	info, err := os.Stat(canonical)
	if err != nil {
		return "", "", fmt.Errorf("path is unavailable: %s", filepath.ToSlash(clean))
	}
	if requireDirectory && !info.IsDir() {
		return "", "", fmt.Errorf("path is not a directory: %s", filepath.ToSlash(clean))
	}
	if !requireDirectory && !info.Mode().IsRegular() {
		return "", "", fmt.Errorf("path is not a regular file: %s", filepath.ToSlash(clean))
	}
	return canonical, filepath.ToSlash(clean), nil
}

func truncateUTF8(value string, limit int) (string, bool) {
	if limit < 0 || len([]byte(value)) <= limit {
		return value, false
	}
	if limit == 0 {
		return "", true
	}
	marker := []byte("\n\n[Result truncated by the frozen tool policy.]")
	contentLimit := limit
	if limit > len(marker) {
		contentLimit -= len(marker)
	} else {
		marker = nil
	}
	data := []byte(value)
	data = data[:contentLimit]
	for len(data) > 0 && !utf8.Valid(data) {
		data = data[:len(data)-1]
	}
	return string(append(data, marker...)), true
}

func linePreview(value string, max int) string {
	value = strings.ReplaceAll(value, "\t", "    ")
	if len([]byte(value)) <= max {
		return value
	}
	trimmed, _ := truncateUTF8(value, max)
	return trimmed + "…"
}

func looksBinary(data []byte) bool {
	probe := data
	if len(probe) > 8192 {
		probe = probe[:8192]
	}
	return bytes.IndexByte(probe, 0) >= 0 || !utf8.Valid(probe)
}

func scanTextFile(path string, maxBytes int64, visit func(line int, text string) bool) (bool, error) {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxBytes {
		return false, err
	}
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	probe := make([]byte, 8192)
	n, readErr := file.Read(probe)
	if readErr != nil && readErr != io.EOF {
		return false, readErr
	}
	if looksBinary(probe[:n]) {
		return false, nil
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	line := 0
	for scanner.Scan() {
		if !utf8.Valid(scanner.Bytes()) || bytes.IndexByte(scanner.Bytes(), 0) >= 0 {
			return false, nil
		}
		line++
		if !visit(line, scanner.Text()) {
			return true, nil
		}
	}
	return true, scanner.Err()
}
