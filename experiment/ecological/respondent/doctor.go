package respondent

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
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

type DoctorResult struct {
	SchemaVersion int            `json:"schema_version"`
	TaskID        string         `json:"task_id"`
	Snapshot      map[string]any `json:"snapshot"`
	Policy        map[string]any `json:"policy"`
	Extension     map[string]any `json:"extension"`
	Runtime       map[string]any `json:"runtime"`
	Isolation     map[string]any `json:"isolation"`
	RPC           map[string]any `json:"rpc"`
	Artifacts     map[string]any `json:"artifacts"`
}

func Doctor(root, snapshotAttempt, configPath, artifacts string) (DoctorResult, error) {
	config, configData, err := LoadConfig(root, configPath)
	if err != nil {
		return DoctorResult{}, err
	}
	if !filepath.IsAbs(snapshotAttempt) {
		snapshotAttempt = filepath.Join(root, snapshotAttempt)
	}
	verified, err := mediator.VerifySnapshotAttempt(snapshotAttempt, mediator.AttemptRequirements{
		TaskID: config.TaskID, ControllerCommit: config.SnapshotControllerCommit,
		SourceTreeSHA256: config.SourceTreeSHA256,
	})
	if err != nil {
		return DoctorResult{}, err
	}
	policyPath := rooted(root, config.ToolPolicy)
	policy, _, err := mediator.LoadPolicy(policyPath)
	if err != nil {
		return DoctorResult{}, err
	}
	if policy.TaskID != config.TaskID {
		return DoctorResult{}, fmt.Errorf("policy task does not match respondent config")
	}
	for _, tool := range config.ActiveTools {
		if tool == "submit_ecological_review" {
			continue
		}
		found := false
		for _, allowed := range policy.AllowedTools {
			if tool == allowed {
				found = true
			}
		}
		if !found {
			return DoctorResult{}, fmt.Errorf("active tool %q is not in mediator policy", tool)
		}
	}
	extensionPath := rooted(root, config.PiExtension)
	extensionDigest, err := snapshot.FileSHA256(extensionPath)
	if err != nil {
		return DoctorResult{}, err
	}
	policyFileDigest, err := snapshot.FileSHA256(policyPath)
	if err != nil {
		return DoctorResult{}, err
	}
	policyCanonicalDigest, err := policy.Digest()
	if err != nil {
		return DoctorResult{}, err
	}
	configDigest := shaData(configData)
	if err := os.MkdirAll(artifacts, 0o700); err != nil {
		return DoctorResult{}, err
	}
	transcriptPath := filepath.Join(artifacts, "mediator-transcript.jsonl")
	transcript, err := os.OpenFile(transcriptPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return DoctorResult{}, err
	}
	defer transcript.Close()
	testRunner, err := mediator.NewGortexTestRunner(verified.Attempt, filepath.Join(artifacts, "focused-tests"))
	if err != nil {
		return DoctorResult{}, err
	}
	session, err := mediator.NewSession(verified.Source, policy, testRunner, transcript)
	if err != nil {
		return DoctorResult{}, err
	}
	runtime, err := ResolvePiRuntime()
	if err != nil {
		return DoctorResult{}, err
	}
	sandbox, err := PrepareChildSandbox(artifacts, runtime, extensionPath)
	if err != nil {
		return DoctorResult{}, err
	}
	sandbox.Environment = append(sandbox.Environment, "COUNCIL_ECOLOGICAL_DOCTOR=1")
	denied := []string{
		filepath.Join(verified.Source, "LICENSE.md"),
		filepath.Join(verified.Controller, "provenance.json"),
		filepath.Join(root, "README.md"),
		filepath.Join(root, "experiment", "ecological", "candidates.json"),
		filepath.Join(root, "experiment", "ecological", "evidence", "eco-gortex-unicode-tokenizer.md"),
	}
	for _, path := range denied {
		if data, err := os.ReadFile(path); err != nil || len(data) == 0 {
			return DoctorResult{}, fmt.Errorf("doctor denial sentinel unavailable: %s", path)
		}
	}
	probe, err := ProbeChildBoundary(sandbox, runtime, extensionPath, denied)
	if writeErr := writeJSON(filepath.Join(artifacts, "pi-boundary-probe.json"), probe); writeErr != nil {
		return DoctorResult{}, writeErr
	}
	if err != nil {
		return DoctorResult{}, err
	}

	rpc, err := runRPCHealth(config, runtime, sandbox, extensionPath, session, artifacts)
	if err != nil {
		return DoctorResult{}, err
	}
	if err := transcript.Sync(); err != nil {
		return DoctorResult{}, err
	}
	transcriptDigest, err := snapshot.FileSHA256(transcriptPath)
	if err != nil {
		return DoctorResult{}, err
	}
	profileDigest, err := snapshot.FileSHA256(sandbox.Profile)
	if err != nil {
		return DoctorResult{}, err
	}
	probeProfileDigest, err := snapshot.FileSHA256(filepath.Join(sandbox.Root, "probe-profile.sb"))
	if err != nil {
		return DoctorResult{}, err
	}
	nodeDigest, err := snapshot.FileSHA256(runtime.Node)
	if err != nil {
		return DoctorResult{}, err
	}
	scriptDigest, err := snapshot.FileSHA256(runtime.Script)
	if err != nil {
		return DoctorResult{}, err
	}
	nodeVersion, err := capture(runtime.Node, "--version")
	if err != nil {
		return DoctorResult{}, err
	}
	piVersion, err := capture(runtime.Node, runtime.Script, "--version")
	if err != nil {
		return DoctorResult{}, err
	}
	return DoctorResult{
		SchemaVersion: 1, TaskID: config.TaskID,
		Snapshot: map[string]any{
			"controller_commit": verified.ControllerCommit, "source_tree_sha256": verified.SourceTreeSHA256,
			"source_entries": verified.SourceEntries, "source_bytes": verified.SourceBytes,
		},
		Policy: map[string]any{
			"path": config.ToolPolicy, "file_sha256": policyFileDigest,
			"canonical_sha256": policyCanonicalDigest, "max_calls": policy.MaxCalls,
		},
		Extension: map[string]any{"path": config.PiExtension, "sha256": extensionDigest, "active_tools": config.ActiveTools},
		Runtime: map[string]any{
			"node_version": strings.TrimSpace(nodeVersion), "node_sha256": nodeDigest,
			"pi_version": strings.TrimSpace(piVersion), "pi_entrypoint_sha256": scriptDigest,
			"model_selected_without_call": config.Provider.Model, "thinking": config.Provider.Thinking,
		},
		Isolation: map[string]any{
			"seatbelt": true, "profile_sha256": profileDigest, "probe_profile_sha256": probeProfileDigest,
			"empty_cwd": true, "ephemeral_home": true, "provider_transport_network_only": true,
			"provider_child_source_read_denied": true, "controller_read_denied": true,
			"evidence_read_denied": true, "council_read_denied_except_literal_extension": true,
			"scratch_write_allowed": true, "builtin_tools_disabled": true,
		},
		RPC: rpc,
		Artifacts: map[string]any{
			"config_sha256": configDigest, "mediator_transcript_sha256": transcriptDigest,
			"pi_events": "pi-events.jsonl", "pi_stderr": "pi-stderr.txt",
		},
	}, nil
}

func runRPCHealth(config Config, runtime PiRuntime, sandbox ChildSandbox, extensionPath string, session *mediator.Session, artifacts string) (map[string]any, error) {
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
	defer requestReader.Close()
	defer responseWriter.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.DoctorTimeoutSeconds)*time.Second)
	defer cancel()
	tools := strings.Join(config.ActiveTools, ",")
	args := []string{
		runtime.Script, "--mode", "rpc", "--no-session", "--no-builtin-tools", "--tools", tools,
		"--no-extensions", "--extension", extensionPath, "--no-skills", "--no-prompt-templates",
		"--no-themes", "--no-context-files", "--no-approve", "--model", config.Provider.Model,
		"--thinking", config.Provider.Thinking,
	}
	commandArgs := append([]string{"-f", sandbox.Profile, "--", runtime.Node}, args...)
	command := exec.CommandContext(ctx, "/usr/bin/sandbox-exec", commandArgs...)
	command.Dir = sandbox.CWD
	command.Env = sandbox.Environment
	command.ExtraFiles = []*os.File{requestWriter, responseReader}
	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderrPath := filepath.Join(artifacts, "pi-stderr.txt")
	stderrFile, err := os.OpenFile(stderrPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return nil, err
	}
	defer stderrFile.Close()
	command.Stderr = stderrFile
	if err := command.Start(); err != nil {
		return nil, err
	}
	requestWriter.Close()
	responseReader.Close()
	serveDone := make(chan error, 1)
	go func() {
		err := session.Serve(requestReader, responseWriter)
		_ = responseWriter.Close()
		serveDone <- err
	}()
	eventsPath := filepath.Join(artifacts, "pi-events.jsonl")
	events, err := os.OpenFile(eventsPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	reader := bufio.NewReader(stdout)
	healthCommand := map[string]any{"id": "health", "type": "prompt", "message": "/ecological-mediator-health"}
	if err := sendRPC(stdin, healthCommand); err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	healthResponse, healthNotify, extensionErrors, err := awaitHealth(reader, events)
	if err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	if !healthResponse || !healthNotify || extensionErrors != 0 {
		_ = command.Process.Kill()
		return nil, fmt.Errorf("Pi ecological extension health check did not complete cleanly")
	}
	if err := sendRPC(stdin, map[string]any{"id": "thinking", "type": "set_thinking_level", "level": config.Provider.Thinking}); err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	moreExtensionErrors, err := awaitCommand(reader, events, "thinking")
	extensionErrors += moreExtensionErrors
	if err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	if err := sendRPC(stdin, map[string]any{"id": "state", "type": "get_state"}); err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	stateModel, stateThinking, moreExtensionErrors, err := awaitState(reader, events)
	extensionErrors += moreExtensionErrors
	if err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	if stateModel != config.Provider.Model || stateThinking != config.Provider.Thinking {
		_ = command.Process.Kill()
		return nil, fmt.Errorf("Pi selected model/thinking %q/%q, expected %q/%q", stateModel, stateThinking, config.Provider.Model, config.Provider.Thinking)
	}
	_ = stdin.Close()
	if err := command.Process.Kill(); err != nil {
		return nil, fmt.Errorf("stop Pi RPC doctor after completed checks: %w", err)
	}
	_, _ = io.Copy(events, reader)
	if err := events.Sync(); err != nil {
		_ = command.Process.Kill()
		return nil, err
	}
	if err := events.Close(); err != nil {
		return nil, err
	}
	waitErr := command.Wait()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, fmt.Errorf("Pi RPC doctor timed out")
	}
	if waitErr == nil {
		return nil, fmt.Errorf("Pi RPC doctor unexpectedly exited before controlled termination")
	}
	if exitError, ok := waitErr.(*exec.ExitError); !ok || exitError.ProcessState.ExitCode() != -1 {
		_ = stderrFile.Sync()
		stderrData, _ := os.ReadFile(stderrPath)
		return nil, fmt.Errorf("Pi RPC doctor ended unexpectedly: %w: %s", waitErr, strings.TrimSpace(string(stderrData)))
	}
	serveErr := <-serveDone
	if serveErr != nil {
		return nil, serveErr
	}
	if extensionErrors != 0 {
		return nil, fmt.Errorf("Pi emitted %d extension errors", extensionErrors)
	}
	if err := stderrFile.Sync(); err != nil {
		return nil, err
	}
	eventData, err := os.ReadFile(eventsPath)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(eventData, []byte(`"type":"agent_start"`)) || bytes.Contains(eventData, []byte(`"type":"turn_start"`)) {
		return nil, fmt.Errorf("Pi RPC doctor unexpectedly started an agent/provider turn")
	}
	eventsDigest, err := snapshot.FileSHA256(eventsPath)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"mode": "rpc", "health_command_handled": true, "inherited_pipes_exercised": true,
		"state_model": stateModel, "state_thinking": stateThinking,
		"thinking_set_via_rpc": true, "controlled_termination_after_checks": "SIGKILL", "extension_errors": extensionErrors,
		"provider_calls": 0, "events_sha256": eventsDigest, "exit_code": -1,
	}, nil
}

func awaitHealth(reader *bufio.Reader, events io.Writer) (bool, bool, int, error) {
	response, notify, errors := false, false, 0
	for !response || !notify {
		value, err := readRPC(reader, events)
		if err != nil {
			return response, notify, errors, err
		}
		if value["type"] == "extension_error" {
			errors++
		}
		if value["type"] == "response" && value["id"] == "health" {
			if value["success"] != true {
				return false, notify, errors, fmt.Errorf("health command failed: %v", value["error"])
			}
			response = true
		}
		if value["type"] == "extension_ui_request" && value["method"] == "notify" {
			message, _ := value["message"].(string)
			if strings.Contains(message, `"ecological_mediator_health":true`) {
				notify = true
			}
		}
	}
	return response, notify, errors, nil
}

func awaitCommand(reader *bufio.Reader, events io.Writer, id string) (int, error) {
	errors := 0
	for {
		value, err := readRPC(reader, events)
		if err != nil {
			return errors, err
		}
		if value["type"] == "extension_error" {
			errors++
		}
		if value["type"] != "response" || value["id"] != id {
			continue
		}
		if value["success"] != true {
			return errors, fmt.Errorf("RPC command %s failed: %v", id, value["error"])
		}
		return errors, nil
	}
}

func awaitState(reader *bufio.Reader, events io.Writer) (string, string, int, error) {
	errors := 0
	for {
		value, err := readRPC(reader, events)
		if err != nil {
			return "", "", errors, err
		}
		if value["type"] == "extension_error" {
			errors++
		}
		if value["type"] != "response" || value["id"] != "state" {
			continue
		}
		if value["success"] != true {
			return "", "", errors, fmt.Errorf("get_state failed: %v", value["error"])
		}
		data, _ := value["data"].(map[string]any)
		model, _ := data["model"].(map[string]any)
		id, _ := model["id"].(string)
		thinking, _ := data["thinkingLevel"].(string)
		return id, thinking, errors, nil
	}
}

func readRPC(reader *bufio.Reader, events io.Writer) (map[string]any, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read Pi RPC: %w", err)
	}
	if _, err := events.Write(line); err != nil {
		return nil, err
	}
	line = bytes.TrimSuffix(line, []byte{'\n'})
	line = bytes.TrimSuffix(line, []byte{'\r'})
	var value map[string]any
	if err := json.Unmarshal(line, &value); err != nil {
		return nil, fmt.Errorf("parse Pi RPC frame: %w", err)
	}
	return value, nil
}

func sendRPC(writer io.Writer, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = writer.Write(append(data, '\n'))
	return err
}

func rooted(root, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, filepath.FromSlash(path))
}

func capture(name string, args ...string) (string, error) {
	command := exec.Command(name, args...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %w: %s", name, err, strings.TrimSpace(stderr.String()))
	}
	return string(output), nil
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o600)
}

func shaData(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}
