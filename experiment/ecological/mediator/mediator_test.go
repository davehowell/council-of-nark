package mediator

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

func testPolicy() Policy {
	return Policy{
		SchemaVersion: 1, TaskID: "test-task",
		AllowedTools: []string{"source_list", "source_read", "source_search", "run_focused_test"},
		MaxCalls:     8, MaxCallsByTool: map[string]int{
			"source_list": 2, "source_read": 3, "source_search": 2, "run_focused_test": 1,
		},
		MaxResultBytes: 2048, MaxTotalResultBytes: 8192, MaxReadLines: 20,
		MaxListDepth: 3, MaxListEntries: 20, MaxSearchMatches: 5,
		MaxSearchPattern: 100, MaxTextFileBytes: 4096,
		TestTargets: []string{"focused"},
	}
}

func testSource(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"README.md":         "alpha\nbeta\ngamma\n",
		"pkg/token.go":      "package pkg\n\nfunc Tokenize(s string) string { return s }\n// ТЕКСТ\n",
		"pkg/token_test.go": "package pkg\n// Tokenize exercises Unicode\n",
	}
	for path, content := range files {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o444); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func request(id, tool, args string) Request {
	return Request{RequestID: id, Tool: tool, Arguments: json.RawMessage(args)}
}

func TestReadListSearchAndTranscript(t *testing.T) {
	root := testSource(t)
	var transcript bytes.Buffer
	session, err := NewSession(root, testPolicy(), TestRunnerFunc(func(target string) (string, map[string]any, error) {
		return "FAIL: controlled panic", map[string]any{"target": target, "network_denied": true}, nil
	}), &transcript)
	if err != nil {
		t.Fatal(err)
	}

	listed := session.Handle(request("1", "source_list", `{"path":"pkg","depth":1}`))
	if !listed.OK || listed.Sequence != 1 || !strings.Contains(listed.Output, "f pkg/token.go") {
		t.Fatalf("unexpected list response: %#v", listed)
	}
	read := session.Handle(request("2", "source_read", `{"path":"pkg/token.go","start_line":3,"line_count":2}`))
	if !read.OK || !strings.Contains(read.Output, "3:func Tokenize") || !strings.Contains(read.Output, "4:// ТЕКСТ") {
		t.Fatalf("unexpected read response: %#v", read)
	}
	searched := session.Handle(request("3", "source_search", `{"pattern":"Tokenize|ТЕКСТ","path":"pkg","include":"*.go"}`))
	if !searched.OK || !strings.Contains(searched.Output, "pkg/token.go:3") || !strings.Contains(searched.Output, "pkg/token_test.go:2") {
		t.Fatalf("unexpected search response: %#v", searched)
	}
	tested := session.Handle(request("4", "run_focused_test", `{"target":"focused"}`))
	if !tested.OK || !strings.Contains(tested.Output, "controlled panic") || tested.Metadata["network_denied"] != true {
		t.Fatalf("unexpected test response: %#v", tested)
	}
	if lines := strings.Count(strings.TrimSpace(transcript.String()), "\n") + 1; lines != 4 {
		t.Fatalf("got %d transcript lines", lines)
	}
}

func TestResultTruncationStaysInsideByteBudget(t *testing.T) {
	value := strings.Repeat("é", 100)
	truncated, wasTruncated := truncateUTF8(value, 80)
	if !wasTruncated || len([]byte(truncated)) > 80 || !strings.Contains(truncated, "Result truncated") {
		t.Fatalf("unexpected truncation: %d bytes %q", len([]byte(truncated)), truncated)
	}
}

func TestPathBoundaryRejectsTraversalGitAndSymlinks(t *testing.T) {
	root := testSource(t)
	outside := filepath.Join(t.TempDir(), "secret.txt")
	if err := os.WriteFile(outside, []byte("secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	policy := testPolicy()
	policy.MaxCallsByTool["source_read"] = 10
	session, err := NewSession(root, policy, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for index, path := range []string{"../secret.txt", ".git/HEAD", "link", "/etc/passwd"} {
		response := session.Handle(request(string(rune('a'+index)), "source_read", `{"path":`+quoted(path)+`}`))
		if response.OK || response.Error == "" {
			t.Fatalf("path %q was not denied: %#v", path, response)
		}
	}
}

func TestSourceRootWithGitIsRejected(t *testing.T) {
	root := testSource(t)
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := NewSession(root, testPolicy(), nil, nil); err == nil {
		t.Fatal("source root containing .git was accepted")
	}
}

func TestBudgetsAndStrictArguments(t *testing.T) {
	policy := testPolicy()
	policy.MaxCallsByTool["source_read"] = 2
	session, err := NewSession(testSource(t), policy, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	bad := session.Handle(request("bad", "source_read", `{"path":"README.md","unknown":true}`))
	if bad.OK || !strings.Contains(bad.Error, "unknown field") {
		t.Fatalf("unknown argument accepted: %#v", bad)
	}
	first := session.Handle(request("first", "source_read", `{"path":"README.md","line_count":1}`))
	if !first.OK {
		t.Fatal(first.Error)
	}
	second := session.Handle(request("second", "source_read", `{"path":"README.md","line_count":1}`))
	if second.OK || !strings.Contains(second.Error, "budget exhausted") {
		t.Fatalf("per-tool budget was not enforced: %#v", second)
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, os.ErrPermission }

func TestTranscriptFailureClosesSession(t *testing.T) {
	session, err := NewSession(testSource(t), testPolicy(), nil, failingWriter{})
	if err != nil {
		t.Fatal(err)
	}
	first := session.Handle(request("first", "source_read", `{"path":"README.md","line_count":1}`))
	if first.OK || !strings.Contains(first.Error, "transcript write failed") {
		t.Fatalf("transcript failure did not fail closed: %#v", first)
	}
	second := session.Handle(request("second", "source_read", `{"path":"README.md","line_count":1}`))
	if second.OK || !strings.Contains(second.Error, "transcript is unavailable") {
		t.Fatalf("session continued after transcript failure: %#v", second)
	}
}

func TestServeCorrelatesFrames(t *testing.T) {
	session, err := NewSession(testSource(t), testPolicy(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	input := strings.NewReader("{bad json}\n" + `{"request_id":"ok","tool":"source_read","arguments":{"path":"README.md","line_count":1}}` + "\n")
	var output bytes.Buffer
	if err := session.Serve(input, &output); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], "invalid request frame") || !strings.Contains(lines[1], `"request_id":"ok"`) {
		t.Fatalf("unexpected protocol output: %s", output.String())
	}
}

func TestVerifySnapshotAttemptDetectsTampering(t *testing.T) {
	attempt := t.TempDir()
	source := filepath.Join(attempt, "source")
	if err := os.Mkdir(source, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "main.go"), []byte("package main\n"), 0o444); err != nil {
		t.Fatal(err)
	}
	controller := filepath.Join(attempt, "controller")
	if err := os.Mkdir(controller, 0o700); err != nil {
		t.Fatal(err)
	}
	manifest, err := snapshot.BuildTreeManifest(source)
	if err != nil {
		t.Fatal(err)
	}
	writeTestJSON(t, filepath.Join(attempt, "status.json"), map[string]any{"task_id": "task", "status": "success", "phase": "complete"})
	writeTestJSON(t, filepath.Join(controller, "source-manifest.json"), manifest)
	writeTestJSON(t, filepath.Join(controller, "provenance.json"), map[string]any{
		"task_id": "task", "controller": map[string]any{"commit": "abc", "tree_dirty": false},
		"validation": map[string]any{"closure_verified_offline": true, "network_probe_denied": true},
		"isolation": map[string]any{
			"source_read": true, "source_write_denied": true, "council_read_denied": true,
			"controller_metadata_read_denied": true, "sibling_task_read_denied": true,
			"evidence_read_denied": true, "git_history_absent": true, "network_denied": true,
		},
	})
	provenanceSHA, _ := snapshot.FileSHA256(filepath.Join(controller, "provenance.json"))
	manifestSHA, _ := snapshot.FileSHA256(filepath.Join(controller, "source-manifest.json"))
	writeTestJSON(t, filepath.Join(controller, "seal.json"), map[string]any{
		"provenance_sha256": provenanceSHA, "source_manifest_sha256": manifestSHA, "source_tree_sha256": manifest.TreeSHA256,
	})
	requirements := AttemptRequirements{TaskID: "task", ControllerCommit: "abc", SourceTreeSHA256: manifest.TreeSHA256}
	verified, err := VerifySnapshotAttempt(attempt, requirements)
	if err != nil || verified.SourceEntries != 1 {
		t.Fatalf("valid attempt rejected: %#v %v", verified, err)
	}
	if err := os.Chmod(filepath.Join(source, "main.go"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "main.go"), []byte("tampered\n"), 0o444); err != nil {
		t.Fatal(err)
	}
	if _, err := VerifySnapshotAttempt(attempt, requirements); err == nil || !strings.Contains(err.Error(), "source tree") {
		t.Fatalf("tampered source was accepted: %v", err)
	}
}

func TestGortexProfileHasNoNetworkAllowanceAndSanitizesPaths(t *testing.T) {
	profile := gortexTestProfile("/private/parent", "/private/go", "/private/mod", "/private/go/bin/go", "/private/scratch")
	if strings.Contains(profile, "network-outbound") || !strings.Contains(profile, "(deny default)") {
		t.Fatalf("unsafe focused-test profile:\n%s", profile)
	}
	runner := &GortexTestRunner{ParentRoot: "/secret/parent", Toolchain: "/secret/toolchain", Modules: "/secret/modules", Ecological: "/secret"}
	output := runner.sanitizeOutput("/secret/parent/a.go /secret/toolchain/root/src/x /secret/modules/gomod/m", "/tmp/scratch")
	if strings.Contains(output, "/secret") || !strings.Contains(output, "<TEST_ROOT>") || !strings.Contains(output, "<FROZEN_GOROOT>") {
		t.Fatalf("paths were not sanitized: %s", output)
	}
}

func TestTrackedPolicyLoadsAndDigests(t *testing.T) {
	path := filepath.Join("..", "tool-policy", "eco-gortex-unicode-tokenizer-v1.json")
	policy, data, err := LoadPolicy(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || policy.TaskID != "eco-gortex-unicode-tokenizer" {
		t.Fatal("tracked policy did not load")
	}
	first, err := policy.Digest()
	if err != nil {
		t.Fatal(err)
	}
	policy.AllowedTools[0], policy.AllowedTools[1] = policy.AllowedTools[1], policy.AllowedTools[0]
	second, err := policy.Digest()
	if err != nil || first != second {
		t.Fatalf("canonical policy digest changed with list order: %s %s %v", first, second, err)
	}
}

func writeTestJSON(t *testing.T, path string, value any) {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
}

func quoted(value string) string {
	data, _ := json.Marshal(value)
	return string(data)
}
