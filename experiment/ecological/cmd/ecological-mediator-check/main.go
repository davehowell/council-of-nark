package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/davehowell/council-of-nark/experiment/ecological/mediator"
	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

const (
	taskID           = "eco-gortex-unicode-tokenizer"
	controllerCommit = "a0004597c5eb74bd4ea0697304301d06e4b8d55f"
	sourceTreeSHA256 = "41a9a14bd447791f254c79a7441d81451dd69ce6e03367dddbc20c4679c2b523"
)

func main() {
	snapshotPath := flag.String("snapshot", "", "completed clean Gortex snapshot attempt")
	policyPath := flag.String("policy", "experiment/ecological/tool-policy/eco-gortex-unicode-tokenizer-v1.json", "tracked tool policy")
	flag.Parse()
	if *snapshotPath == "" || flag.NArg() != 0 {
		fmt.Fprintln(os.Stderr, "usage: ecological-mediator-check --snapshot <attempt> [--policy <path>]")
		os.Exit(2)
	}
	root, err := repositoryRoot()
	if err != nil {
		fail("find repository root", err)
	}
	attempt, err := newAttempt(filepath.Join(root, "experiment", "ecological", "mediator-runs"))
	if err != nil {
		fail("create check attempt", err)
	}
	status := map[string]any{
		"schema_version": 1, "status": "running", "started_at": time.Now().UTC().Format(time.RFC3339Nano),
		"task_id": taskID,
	}
	_ = writeJSON(filepath.Join(attempt, "status.json"), status)
	if err := run(root, attempt, *snapshotPath, *policyPath); err != nil {
		status["status"] = "failed"
		status["error"] = err.Error()
		status["finished_at"] = time.Now().UTC().Format(time.RFC3339Nano)
		_ = writeJSON(filepath.Join(attempt, "status.json"), status)
		printAttempt(root, attempt)
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	summaryDigest, sealErr := snapshot.FileSHA256(filepath.Join(attempt, "summary.json"))
	if sealErr == nil {
		sealErr = writeJSON(filepath.Join(attempt, "seal.json"), map[string]any{
			"schema_version": 1, "summary_sha256": summaryDigest,
			"completed_at": time.Now().UTC().Format(time.RFC3339Nano),
		})
	}
	if sealErr != nil {
		status["status"] = "failed"
		status["error"] = "seal check: " + sealErr.Error()
		status["finished_at"] = time.Now().UTC().Format(time.RFC3339Nano)
		_ = writeJSON(filepath.Join(attempt, "status.json"), status)
		printAttempt(root, attempt)
		fmt.Fprintln(os.Stderr, "error:", sealErr)
		os.Exit(1)
	}
	status["status"] = "success"
	status["finished_at"] = time.Now().UTC().Format(time.RFC3339Nano)
	_ = writeJSON(filepath.Join(attempt, "status.json"), status)
	printAttempt(root, attempt)
}

func run(root, attempt, attemptArgument, policyArgument string) error {
	if !filepath.IsAbs(attemptArgument) {
		attemptArgument = filepath.Join(root, attemptArgument)
	}
	verified, err := mediator.VerifySnapshotAttempt(attemptArgument, mediator.AttemptRequirements{
		TaskID: taskID, ControllerCommit: controllerCommit, SourceTreeSHA256: sourceTreeSHA256,
	})
	if err != nil {
		return err
	}
	if !filepath.IsAbs(policyArgument) {
		policyArgument = filepath.Join(root, policyArgument)
	}
	policy, policyData, err := mediator.LoadPolicy(policyArgument)
	if err != nil {
		return err
	}
	if policy.TaskID != taskID {
		return fmt.Errorf("tool policy task mismatch")
	}
	policyDigest, err := policy.Digest()
	if err != nil {
		return err
	}
	policyFileDigest, err := snapshot.FileSHA256(policyArgument)
	if err != nil {
		return err
	}
	extensionPath := filepath.Join(root, "experiment", "ecological", "pi", "ecological-tools.ts")
	extensionDigest, err := snapshot.FileSHA256(extensionPath)
	if err != nil {
		return err
	}
	mediatorManifest, err := snapshot.BuildTreeManifest(filepath.Join(root, "experiment", "ecological", "mediator"))
	if err != nil {
		return err
	}
	currentCommit, err := gitOutput(root, "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	currentStatus, err := gitOutput(root, "status", "--porcelain=v1")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(attempt, "policy.json"), policyData, 0o600); err != nil {
		return err
	}
	transcript, err := os.OpenFile(filepath.Join(attempt, "transcript.jsonl"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer transcript.Close()
	testRunner, err := mediator.NewGortexTestRunner(verified.Attempt, filepath.Join(attempt, "tests"))
	if err != nil {
		return err
	}
	session, err := mediator.NewSession(verified.Source, policy, testRunner, transcript)
	if err != nil {
		return err
	}
	requests := []mediator.Request{
		makeRequest("list", "source_list", map[string]any{"path": "internal/search/rerank", "depth": 1}),
		makeRequest("read", "source_read", map[string]any{"path": "internal/search/rerank/tokens.go", "start_line": 1, "line_count": 80}),
		makeRequest("search", "source_search", map[string]any{"pattern": `\[\]rune|unicode\.IsUpper`, "path": "internal/search/rerank", "include": "*.go"}),
		makeRequest("test", "run_focused_test", map[string]any{"target": "unicode-tokenizer-regression"}),
		makeRequest("deny-traversal", "source_read", map[string]any{"path": "../controller/provenance.json"}),
		makeRequest("deny-target", "run_focused_test", map[string]any{"target": "arbitrary-shell"}),
	}
	responses := make([]mediator.Response, 0, len(requests))
	for _, request := range requests {
		responses = append(responses, session.Handle(request))
	}
	for index := 0; index < 4; index++ {
		if !responses[index].OK {
			return fmt.Errorf("positive mediator probe %s failed: %s", requests[index].RequestID, responses[index].Error)
		}
	}
	for index := 4; index < len(responses); index++ {
		if responses[index].OK || responses[index].Error == "" {
			return fmt.Errorf("negative mediator probe %s was not denied", requests[index].RequestID)
		}
	}
	if err := transcript.Sync(); err != nil {
		return err
	}
	transcriptDigest, err := snapshot.FileSHA256(filepath.Join(attempt, "transcript.jsonl"))
	if err != nil {
		return err
	}
	testArtifacts, err := snapshot.BuildTreeManifest(filepath.Join(attempt, "tests"))
	if err != nil {
		return err
	}
	return writeJSON(filepath.Join(attempt, "summary.json"), map[string]any{
		"schema_version": 1, "task_id": taskID,
		"controller": map[string]any{
			"commit": strings.TrimSpace(currentCommit), "tree_dirty": strings.TrimSpace(currentStatus) != "",
			"mediator_tree_sha256": mediatorManifest.TreeSHA256,
		},
		"snapshot": map[string]any{
			"controller_commit": verified.ControllerCommit, "source_tree_sha256": verified.SourceTreeSHA256,
			"source_entries": verified.SourceEntries, "source_bytes": verified.SourceBytes,
		},
		"policy_canonical_sha256": policyDigest, "policy_file_sha256": policyFileDigest,
		"extension_sha256":           extensionDigest,
		"transcript_sha256":          transcriptDigest,
		"test_artifacts_tree_sha256": testArtifacts.TreeSHA256, "responses": responses,
		"probes": map[string]any{
			"list": true, "read": true, "search": true, "focused_test": true,
			"traversal_denied": true, "arbitrary_test_target_denied": true,
		},
	})
}

func makeRequest(id, tool string, args any) mediator.Request {
	data, _ := json.Marshal(args)
	return mediator.Request{RequestID: id, Tool: tool, Arguments: data}
}

func gitOutput(root string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func repositoryRoot() (string, error) {
	working, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(working, ".git")); err == nil {
			return working, nil
		}
		parent := filepath.Dir(working)
		if parent == working {
			return "", fmt.Errorf("not inside repository")
		}
		working = parent
	}
}

func newAttempt(parent string) (string, error) {
	if err := os.MkdirAll(parent, 0o700); err != nil {
		return "", err
	}
	var suffix [4]byte
	if _, err := rand.Read(suffix[:]); err != nil {
		return "", err
	}
	path := filepath.Join(parent, time.Now().UTC().Format("20060102T150405Z")+"-gortex-mediator-check-"+hex.EncodeToString(suffix[:]))
	if err := os.Mkdir(path, 0o700); err != nil {
		return "", err
	}
	return path, nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func printAttempt(root, attempt string) {
	if relative, err := filepath.Rel(root, attempt); err == nil {
		fmt.Println(filepath.ToSlash(relative))
		return
	}
	fmt.Println(attempt)
}

func fail(label string, err error) {
	fmt.Fprintf(os.Stderr, "error: %s: %v\n", label, err)
	os.Exit(1)
}
