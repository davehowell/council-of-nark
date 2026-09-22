package respondent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/davehowell/council-of-nark/experiment/ecological/mediator"
	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

type RunResult struct {
	SchemaVersion int               `json:"schema_version"`
	TaskID        string            `json:"task_id"`
	Outcome       string            `json:"outcome"`
	Adapter       map[string]any    `json:"adapter"`
	Snapshot      map[string]any    `json:"snapshot"`
	Controller    map[string]any    `json:"controller"`
	Events        EventSummary      `json:"events"`
	Audit         AuditSummary      `json:"provider_audit"`
	Mediator      MediatorSummary   `json:"mediator"`
	Budgets       map[string]any    `json:"budgets"`
	Validation    map[string]any    `json:"validation"`
	Artifacts     map[string]string `json:"artifacts"`
}

type MediatorSummary struct {
	Calls       int            `json:"calls"`
	DeniedCalls int            `json:"denied_calls"`
	ResultBytes int            `json:"result_bytes"`
	ByTool      map[string]int `json:"by_tool"`
}

type preparedRun struct {
	Config        RunConfig
	Verified      mediator.VerifiedAttempt
	Policy        mediator.Policy
	ExtensionPath string
	SystemPrompt  string
	Brief         string
	Controller    map[string]any
}

func ExecuteRun(root, snapshotAttempt, configPath, artifacts string) (RunResult, error) {
	prepared, configData, err := prepareRun(root, snapshotAttempt, configPath)
	if err != nil {
		return RunResult{}, err
	}
	request, err := buildRequest(root, prepared, configPath, configData)
	if err != nil {
		return RunResult{}, err
	}
	if err := writeJSON(filepath.Join(artifacts, "request.json"), request); err != nil {
		return RunResult{}, err
	}
	transcriptPath := filepath.Join(artifacts, "mediator-transcript.jsonl")
	transcript, err := os.OpenFile(transcriptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return RunResult{}, err
	}
	testRunner, err := mediator.NewGortexTestRunner(prepared.Verified.Attempt, filepath.Join(artifacts, "focused-tests"))
	if err != nil {
		_ = transcript.Close()
		return RunResult{}, err
	}
	session, err := mediator.NewSession(prepared.Verified.Source, prepared.Policy, testRunner, transcript)
	if err != nil {
		_ = transcript.Close()
		return RunResult{}, err
	}
	started := time.Now()
	var adapter map[string]any
	if prepared.Config.Adapter == "mock" {
		adapter, err = executeMock(prepared, session, artifacts)
	} else {
		adapter, err = executePi(root, prepared, session, artifacts)
	}
	if syncErr := transcript.Sync(); err == nil && syncErr != nil {
		err = syncErr
	}
	if closeErr := transcript.Close(); err == nil && closeErr != nil {
		err = closeErr
	}
	if err != nil {
		return RunResult{}, err
	}

	eventsPath := filepath.Join(artifacts, "pi-events.jsonl")
	auditPath := filepath.Join(artifacts, "provider-audit.jsonl")
	events, eventErr := AnalyzeEvents(eventsPath, prepared.Config.Budgets.MaxFinalBytes)
	audit, auditErr := AnalyzeAudit(auditPath)
	mediatorSummary, mediatorErr := analyzeMediatorTranscript(transcriptPath)
	validationErrors := []string{}
	for _, candidate := range []error{eventErr, auditErr, mediatorErr} {
		if candidate != nil {
			validationErrors = append(validationErrors, candidate.Error())
		}
	}
	if audit.ProviderRequests < 1 {
		validationErrors = append(validationErrors, "no exact provider request was captured")
	}
	if audit.SessionStarts != 1 || audit.FinalSubmissions != 1 || audit.ToolsAfterSubmission != 0 {
		validationErrors = append(validationErrors, "provider audit does not contain one clean final-submission lifecycle")
	}
	if audit.ProviderRequests < events.AssistantMessages || audit.AssistantUsage != events.AssistantMessages {
		validationErrors = append(validationErrors, fmt.Sprintf("provider request/usage/message counts are inconsistent: %d/%d/%d", audit.ProviderRequests, audit.AssistantUsage, events.AssistantMessages))
	}
	if audit.ProviderResponses > audit.ProviderRequests {
		validationErrors = append(validationErrors, fmt.Sprintf("provider response count exceeds requests: %d/%d", audit.ProviderResponses, audit.ProviderRequests))
	}
	if audit.ProviderRequests > prepared.Config.Budgets.MaxProviderTurns {
		validationErrors = append(validationErrors, "provider-turn cap was exceeded")
	}
	if events.Usage.Total > prepared.Config.Budgets.MaxTotalTokens && events.AssistantAfterSubmission {
		validationErrors = append(validationErrors, "provider continued after a token-budget overrun")
	}
	valid := len(validationErrors) == 0
	if valid {
		if err := writeJSON(filepath.Join(artifacts, "final-submission.json"), events.Review); err != nil {
			return RunResult{}, err
		}
	}
	overrun := events.Usage.Total - prepared.Config.Budgets.MaxTotalTokens
	if overrun < 0 {
		overrun = 0
	}
	outcome := "valid"
	if !valid {
		outcome = "malformed"
	}
	result := RunResult{
		SchemaVersion: 1, TaskID: prepared.Config.TaskID, Outcome: outcome, Adapter: adapter,
		Snapshot: map[string]any{
			"controller_commit":  prepared.Verified.ControllerCommit,
			"source_tree_sha256": prepared.Verified.SourceTreeSHA256,
			"source_entries":     prepared.Verified.SourceEntries, "source_bytes": prepared.Verified.SourceBytes,
		},
		Controller: prepared.Controller, Events: events, Audit: audit, Mediator: mediatorSummary,
		Budgets: map[string]any{
			"limits": prepared.Config.Budgets, "observed_total_tokens": events.Usage.Total,
			"unavoidable_last_response_overrun_tokens": overrun,
			"elapsed_seconds":                          time.Since(started).Seconds(),
		},
		Validation: map[string]any{"valid": valid, "errors": validationErrors},
		Artifacts: map[string]string{
			"request": "request.json", "events": "pi-events.jsonl", "provider_audit": "provider-audit.jsonl",
			"mediator_transcript": "mediator-transcript.jsonl", "final_submission": optionalArtifact(valid, "final-submission.json"),
		},
	}
	if err := writeJSON(filepath.Join(artifacts, "result.json"), result); err != nil {
		return RunResult{}, err
	}
	return result, nil
}

func prepareRun(root, snapshotAttempt, configPath string) (preparedRun, []byte, error) {
	config, configData, err := LoadRunConfig(root, configPath)
	if err != nil {
		return preparedRun{}, nil, err
	}
	if !filepath.IsAbs(snapshotAttempt) {
		snapshotAttempt = filepath.Join(root, snapshotAttempt)
	}
	verified, err := mediator.VerifySnapshotAttempt(snapshotAttempt, mediator.AttemptRequirements{
		TaskID: config.TaskID, ControllerCommit: config.SnapshotControllerCommit, SourceTreeSHA256: config.SourceTreeSHA256,
	})
	if err != nil {
		return preparedRun{}, nil, err
	}
	policy, _, err := mediator.LoadPolicy(rooted(root, config.ToolPolicy))
	if err != nil {
		return preparedRun{}, nil, err
	}
	if policy.TaskID != config.TaskID {
		return preparedRun{}, nil, fmt.Errorf("tool policy task does not match run config")
	}
	for _, tool := range config.ActiveTools {
		if tool == "submit_ecological_review" {
			continue
		}
		found := false
		for _, allowed := range policy.AllowedTools {
			found = found || tool == allowed
		}
		if !found {
			return preparedRun{}, nil, fmt.Errorf("active tool %q is not in the mediator policy", tool)
		}
	}
	extensionPath := rooted(root, config.PiExtension)
	if _, err := os.Stat(extensionPath); err != nil {
		return preparedRun{}, nil, err
	}
	systemPrompt, err := readPrompt(rooted(root, config.SystemPrompt), "system prompt")
	if err != nil {
		return preparedRun{}, nil, err
	}
	brief, err := readPrompt(rooted(root, config.Brief), "brief")
	if err != nil {
		return preparedRun{}, nil, err
	}
	controller, err := currentControllerState(root)
	if err != nil {
		return preparedRun{}, nil, err
	}
	if config.Adapter == "pi" && controller["tree_dirty"] == true {
		return preparedRun{}, nil, fmt.Errorf("live ecological runs require a clean committed controller")
	}
	return preparedRun{
		Config: config, Verified: verified, Policy: policy, ExtensionPath: extensionPath,
		SystemPrompt: systemPrompt, Brief: brief, Controller: controller,
	}, configData, nil
}

func buildRequest(root string, prepared preparedRun, configPath string, configData []byte) (map[string]any, error) {
	files := map[string]string{
		"config": rooted(root, configPath), "tool_policy": rooted(root, prepared.Config.ToolPolicy),
		"pi_extension": prepared.ExtensionPath, "system_prompt": rooted(root, prepared.Config.SystemPrompt),
		"brief": rooted(root, prepared.Config.Brief),
	}
	digests := map[string]string{"config": shaData(configData)}
	for label, path := range files {
		if label == "config" {
			continue
		}
		digest, err := snapshot.FileSHA256(path)
		if err != nil {
			return nil, err
		}
		digests[label] = digest
	}
	policyDigest, err := prepared.Policy.Digest()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"schema_version": 1, "task_id": prepared.Config.TaskID, "adapter": prepared.Config.Adapter,
		"controller": prepared.Controller,
		"snapshot": map[string]any{
			"controller_commit": prepared.Verified.ControllerCommit, "source_tree_sha256": prepared.Verified.SourceTreeSHA256,
		},
		"provider": map[string]any{
			"model": prepared.Config.Provider.Model, "thinking": prepared.Config.Provider.Thinking,
			"decoding_capture": "the complete serialized provider payload for every turn is retained in provider-audit.jsonl",
		},
		"prompts": map[string]any{
			"system": map[string]any{"path": prepared.Config.SystemPrompt, "sha256": digests["system_prompt"], "text": prepared.SystemPrompt},
			"brief":  map[string]any{"path": prepared.Config.Brief, "sha256": digests["brief"], "text": prepared.Brief},
		},
		"tools": map[string]any{
			"active": prepared.Config.ActiveTools, "policy_path": prepared.Config.ToolPolicy,
			"policy_file_sha256": digests["tool_policy"], "policy_canonical_sha256": policyDigest,
			"extension_path": prepared.Config.PiExtension, "extension_sha256": digests["pi_extension"],
			"schema_capture": "serialized tool schemas are retained in each provider_request payload",
		},
		"budgets": prepared.Config.Budgets,
		"isolation": map[string]any{
			"source_visible_only_through_mediator": true, "provider_transport_network_only": true,
			"builtin_tools_disabled": true, "session_disabled": true, "context_discovery_disabled": true,
		},
		"input_digests": digests,
	}, nil
}

func executeMock(prepared preparedRun, session *mediator.Session, artifacts string) (map[string]any, error) {
	probes := []struct {
		name        string
		request     mediator.Request
		mustSucceed bool
	}{
		{"list", mediator.Request{RequestID: "mock-list", Tool: "source_list", Arguments: json.RawMessage(`{"path":"internal/search/rerank","depth":1}`)}, true},
		{"read", mediator.Request{RequestID: "mock-read", Tool: "source_read", Arguments: json.RawMessage(`{"path":"README.md","start_line":1,"line_count":5}`)}, true},
		{"search", mediator.Request{RequestID: "mock-search", Tool: "source_search", Arguments: json.RawMessage(`{"pattern":"token","path":"internal/search/rerank","include":"*.go"}`)}, true},
		{"focused-test", mediator.Request{RequestID: "mock-test", Tool: "run_focused_test", Arguments: json.RawMessage(`{"target":"unicode-tokenizer-regression"}`)}, true},
		{"git-denial", mediator.Request{RequestID: "mock-git", Tool: "source_read", Arguments: json.RawMessage(`{"path":".git/HEAD","line_count":1}`)}, false},
		{"traversal-denial", mediator.Request{RequestID: "mock-traversal", Tool: "source_read", Arguments: json.RawMessage(`{"path":"../controller/provenance.json","line_count":1}`)}, false},
		{"arbitrary-test-denial", mediator.Request{RequestID: "mock-arbitrary-test", Tool: "run_focused_test", Arguments: json.RawMessage(`{"target":"arbitrary-shell-target"}`)}, false},
	}
	probeLog := []map[string]any{}
	for _, probe := range probes {
		response := session.Handle(probe.request)
		probeLog = append(probeLog, map[string]any{"name": probe.name, "must_succeed": probe.mustSucceed, "response": response})
		if response.OK != probe.mustSucceed {
			return nil, fmt.Errorf("mock boundary probe %s had unexpected result", probe.name)
		}
	}
	if err := writeJSON(filepath.Join(artifacts, "mock-boundary-probes.json"), probeLog); err != nil {
		return nil, err
	}
	review := EcologicalReview{
		ReviewSummary: "Deterministic lifecycle fixture; not a respondent finding.",
		Findings: []EcologicalFinding{{
			Title: "Mock lifecycle finding", Locations: []SourceLocation{{Path: "README.md", LineStart: 1, LineEnd: 1}},
			Claim:              "This claim exists only to exercise structured submission validation.",
			Mechanism:          "The mock adapter emits a fixed event stream without a provider call.",
			Consequence:        "It verifies plumbing and must never be scored as ecological evidence.",
			Correction:         "Do not treat mock output as a respondent result.",
			AllocationAndScope: "No source change or runtime allocation is proposed.",
			RegressionTests:    []RegressionCase{{Case: "Run the mock lifecycle", Expected: "One valid sealed submission", WhyDiscriminating: "It covers request capture, mediation, validation, and sealing."}},
			Evidence:           []string{"controller-generated deterministic fixture"}, Confidence: "high",
		}},
		Uncertainties: []string{"The mock does not test provider behavior."},
	}
	usage := map[string]any{
		"input": prepared.Config.Mock.InputTokens, "output": prepared.Config.Mock.OutputTokens,
		"cacheRead": 0, "cacheWrite": 0, "totalTokens": prepared.Config.Mock.InputTokens + prepared.Config.Mock.OutputTokens,
		"cost": map[string]any{"input": 0, "output": 0, "cacheRead": 0, "cacheWrite": 0, "total": 0},
	}
	events := []any{
		map[string]any{"type": "session", "version": 3, "id": "deterministic-mock", "cwd": "<EMPTY_CWD>"},
		map[string]any{"type": "agent_start"},
		map[string]any{"type": "turn_start"},
		map[string]any{"type": "message_end", "message": map[string]any{"role": "assistant", "content": []any{map[string]any{"type": "toolCall", "id": "mock-submit", "name": "submit_ecological_review", "arguments": review}}, "provider": "mock", "model": prepared.Config.Provider.Model, "usage": usage, "stopReason": "toolUse"}},
		map[string]any{"type": "tool_execution_start", "toolCallId": "mock-submit", "toolName": "submit_ecological_review", "args": review},
		map[string]any{"type": "tool_execution_end", "toolCallId": "mock-submit", "toolName": "submit_ecological_review", "result": map[string]any{"content": []any{map[string]any{"type": "text", "text": "Submitted 1 finding(s)."}}, "details": review}, "isError": false},
		map[string]any{"type": "turn_end", "message": map[string]any{"role": "assistant"}, "toolResults": []any{}},
		map[string]any{"type": "agent_end", "messages": []any{}},
		map[string]any{"type": "agent_settled"},
	}
	if err := writeJSONLines(filepath.Join(artifacts, "pi-events.jsonl"), events); err != nil {
		return nil, err
	}
	audit := []any{
		map[string]any{"schema_version": 1, "sequence": 1, "kind": "session_start", "active_tools": prepared.Config.ActiveTools},
		map[string]any{"schema_version": 1, "sequence": 2, "kind": "provider_request", "provider_turn": 1, "payload": map[string]any{"adapter": "deterministic-mock", "system": prepared.SystemPrompt, "brief": prepared.Brief, "tools": prepared.Config.ActiveTools}},
		map[string]any{"schema_version": 1, "sequence": 3, "kind": "provider_response", "provider_turn": 1, "status": 200, "headers": map[string]any{}},
		map[string]any{"schema_version": 1, "sequence": 4, "kind": "assistant_usage", "provider_turn": 1, "usage": usage, "observed_tokens": prepared.Config.Mock.InputTokens + prepared.Config.Mock.OutputTokens},
		map[string]any{"schema_version": 1, "sequence": 5, "kind": "final_submission", "tool_call_id": "mock-submit", "details": review},
	}
	if err := writeJSONLines(filepath.Join(artifacts, "provider-audit.jsonl"), audit); err != nil {
		return nil, err
	}
	return map[string]any{"kind": "deterministic-mock", "provider_calls": 0, "fixture_version": 1}, nil
}

func executePi(root string, prepared preparedRun, session *mediator.Session, artifacts string) (map[string]any, error) {
	runtime, err := ResolvePiRuntime()
	if err != nil {
		return nil, err
	}
	sandbox, err := PrepareChildSandbox(artifacts, runtime, prepared.ExtensionPath)
	if err != nil {
		return nil, err
	}
	defer removeSandboxSecrets(sandbox)
	sandbox.Environment = append(sandbox.Environment,
		"COUNCIL_ECOLOGICAL_AUDIT=1",
		fmt.Sprintf("COUNCIL_ECOLOGICAL_MAX_PROVIDER_TURNS=%d", prepared.Config.Budgets.MaxProviderTurns),
		fmt.Sprintf("COUNCIL_ECOLOGICAL_MAX_TOTAL_TOKENS=%d", prepared.Config.Budgets.MaxTotalTokens),
	)
	denied := []string{
		filepath.Join(prepared.Verified.Source, "LICENSE.md"), filepath.Join(prepared.Verified.Controller, "provenance.json"),
		filepath.Join(root, "README.md"), filepath.Join(root, "experiment", "ecological", "candidates.json"),
		filepath.Join(root, "experiment", "ecological", "evidence", "eco-gortex-unicode-tokenizer.md"),
	}
	probe, err := ProbeChildBoundary(sandbox, runtime, prepared.ExtensionPath, denied)
	if writeErr := writeJSON(filepath.Join(artifacts, "pi-boundary-probe.json"), probe); writeErr != nil {
		return nil, writeErr
	}
	if err != nil {
		return nil, err
	}
	executableProbes, err := probeForbiddenExecutables(sandbox)
	if writeErr := writeJSON(filepath.Join(artifacts, "pi-executable-denials.json"), executableProbes); writeErr != nil {
		return nil, writeErr
	}
	if err != nil {
		return nil, err
	}

	requestReader, requestWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	responseReader, responseWriter, err := os.Pipe()
	if err != nil {
		requestReader.Close()
		requestWriter.Close()
		return nil, err
	}
	auditReader, auditWriter, err := os.Pipe()
	if err != nil {
		requestReader.Close()
		requestWriter.Close()
		responseReader.Close()
		responseWriter.Close()
		return nil, err
	}
	defer requestReader.Close()
	defer responseWriter.Close()
	defer auditReader.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(prepared.Config.Budgets.TimeoutSeconds)*time.Second)
	defer cancel()
	tools := strings.Join(prepared.Config.ActiveTools, ",")
	args := []string{
		runtime.Script, "--mode", "json", "--no-session", "--no-builtin-tools", "--tools", tools,
		"--no-extensions", "--extension", prepared.ExtensionPath, "--no-skills", "--no-prompt-templates",
		"--no-themes", "--no-context-files", "--no-approve", "--model", prepared.Config.Provider.Model,
		"--thinking", prepared.Config.Provider.Thinking, "--system-prompt", prepared.SystemPrompt, "--", prepared.Brief,
	}
	commandArgs := append([]string{"-f", sandbox.Profile, "--", runtime.Node}, args...)
	command := exec.CommandContext(ctx, "/usr/bin/sandbox-exec", commandArgs...)
	command.Dir = sandbox.CWD
	command.Env = sandbox.Environment
	command.ExtraFiles = []*os.File{requestWriter, responseReader, auditWriter}
	eventsPath := filepath.Join(artifacts, "pi-events.jsonl")
	events, err := os.OpenFile(eventsPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, err
	}
	defer events.Close()
	stderrPath := filepath.Join(artifacts, "pi-stderr.txt")
	stderr, err := os.OpenFile(stderrPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, err
	}
	defer stderr.Close()
	auditPath := filepath.Join(artifacts, "provider-audit.jsonl")
	audit, err := os.OpenFile(auditPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, err
	}
	defer audit.Close()
	command.Stdout = events
	command.Stderr = stderr
	started := time.Now()
	if err := command.Start(); err != nil {
		return nil, err
	}
	requestWriter.Close()
	responseReader.Close()
	auditWriter.Close()
	serveDone := make(chan error, 1)
	go func() {
		err := session.Serve(requestReader, responseWriter)
		_ = responseWriter.Close()
		serveDone <- err
	}()
	auditDone := make(chan error, 1)
	go func() {
		_, err := io.Copy(audit, auditReader)
		auditDone <- err
	}()
	waitErr := command.Wait()
	serveErr := <-serveDone
	auditErr := <-auditDone
	for _, file := range []*os.File{events, stderr, audit} {
		if err := file.Sync(); err != nil {
			return nil, err
		}
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("ecological Pi run timed out after %d seconds", prepared.Config.Budgets.TimeoutSeconds)
	}
	if serveErr != nil {
		return nil, serveErr
	}
	if auditErr != nil {
		return nil, auditErr
	}
	if waitErr != nil {
		stderrData, _ := os.ReadFile(stderrPath)
		return nil, fmt.Errorf("Pi respondent failed: %w: %s", waitErr, strings.TrimSpace(string(stderrData)))
	}
	nodeDigest, _ := snapshot.FileSHA256(runtime.Node)
	scriptDigest, _ := snapshot.FileSHA256(runtime.Script)
	return map[string]any{
		"kind": "pi-json", "exit_code": 0, "duration_seconds": time.Since(started).Seconds(),
		"node_sha256": nodeDigest, "pi_entrypoint_sha256": scriptDigest,
	}, nil
}

func analyzeMediatorTranscript(path string) (MediatorSummary, error) {
	file, err := os.Open(path)
	if err != nil {
		return MediatorSummary{}, err
	}
	defer file.Close()
	summary := MediatorSummary{ByTool: map[string]int{}}
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for scanner.Scan() {
		var entry mediator.TranscriptEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return summary, err
		}
		summary.Calls++
		summary.ByTool[entry.Request.Tool]++
		if !entry.Response.OK {
			summary.DeniedCalls++
		}
		if entry.Response.Metadata != nil {
			summary.ResultBytes += intNumber(entry.Response.Metadata["result_bytes"])
		}
	}
	return summary, scanner.Err()
}

func currentControllerState(root string) (map[string]any, error) {
	commit, err := gitCapture(root, "rev-parse", "HEAD")
	if err != nil {
		return nil, err
	}
	status, err := gitCapture(root, "status", "--porcelain=v1")
	if err != nil {
		return nil, err
	}
	return map[string]any{"commit": strings.TrimSpace(commit), "tree_dirty": strings.TrimSpace(status) != ""}, nil
}

func gitCapture(root string, args ...string) (string, error) {
	command := exec.Command("git", args...)
	command.Dir = root
	output, err := command.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

func readPrompt(path, label string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	text := string(data)
	if strings.TrimSpace(text) == "" || len(data) > 128*1024 {
		return "", fmt.Errorf("%s is empty or too large", label)
	}
	return text, nil
}

func probeForbiddenExecutables(sandbox ChildSandbox) ([]map[string]any, error) {
	probes := []struct {
		name string
		args []string
	}{
		{"shell", []string{"/bin/sh", "-c", "true"}},
		{"git", []string{"/usr/bin/git", "--version"}},
	}
	results := []map[string]any{}
	for _, probe := range probes {
		args := append([]string{"-f", sandbox.Profile, "--"}, probe.args...)
		command := exec.Command("/usr/bin/sandbox-exec", args...)
		command.Dir = sandbox.CWD
		command.Env = sandbox.Environment
		var stdout, stderr bytes.Buffer
		command.Stdout = &stdout
		command.Stderr = &stderr
		err := command.Run()
		exitCode := 0
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else if err != nil {
			exitCode = -1
		}
		results = append(results, map[string]any{"name": probe.name, "exit_code": exitCode, "stdout": stdout.String(), "stderr": stderr.String(), "denied": err != nil})
		if err == nil {
			return results, fmt.Errorf("Pi child profile allowed %s execution", probe.name)
		}
	}
	return results, nil
}

func removeSandboxSecrets(sandbox ChildSandbox) {
	_ = os.RemoveAll(sandbox.Home)
	_ = os.RemoveAll(sandbox.Temp)
	_ = os.RemoveAll(sandbox.CWD)
}

func writeJSONLines(path string, values []any) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	for _, value := range values {
		if err := encoder.Encode(value); err != nil {
			_ = file.Close()
			return err
		}
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func optionalArtifact(condition bool, name string) string {
	if condition {
		return name
	}
	return ""
}
