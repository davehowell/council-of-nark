package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

var worktreeLock sync.Mutex

func (h *Harness) withWorktree(run, callID, commit string, fn func(string) error) error {
	path := filepath.Join(h.Experiment, "worktrees", filepath.Base(run), callID)
	worktreeLock.Lock()
	_ = exec.Command("git", "worktree", "remove", "--force", path).Run()
	_ = os.RemoveAll(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		worktreeLock.Unlock()
		return err
	}
	cmd := exec.Command("git", "worktree", "add", "--detach", "--quiet", path, commit)
	cmd.Dir = h.Root
	out, err := cmd.CombinedOutput()
	worktreeLock.Unlock()
	if err != nil {
		return fmt.Errorf("create detached worktree: %w: %s", err, strings.TrimSpace(string(out)))
	}
	defer func() {
		worktreeLock.Lock()
		remove := exec.Command("git", "worktree", "remove", "--force", path)
		remove.Dir = h.Root
		_ = remove.Run()
		_ = os.RemoveAll(path)
		prune := exec.Command("git", "worktree", "prune")
		prune.Dir = h.Root
		_ = prune.Run()
		worktreeLock.Unlock()
	}()
	actual, err := commandOutput(path, "git", "rev-parse", "HEAD")
	if err != nil {
		return err
	}
	if strings.TrimSpace(string(actual)) != commit {
		return fmt.Errorf("detached worktree is not at frozen commit")
	}
	return fn(path)
}
func (h *Harness) verifyManifest(freeze Freeze) error {
	commit, err := h.sourceCommit()
	if err != nil {
		return err
	}
	if commit != freeze.SourceCommit {
		return fmt.Errorf("HEAD differs from frozen source commit; check out %s before running", freeze.SourceCommit)
	}
	for _, asset := range freeze.Assets {
		data, err := h.gitBlob(freeze.SourceCommit, asset.Path)
		if err != nil {
			return err
		}
		if shaBytes(data) != asset.SHA256 {
			return fmt.Errorf("frozen asset digest mismatch: %s", asset.Path)
		}
	}
	for _, runtimeRow := range freeze.ExternalRuntime {
		for _, raw := range sliceValue(runtimeRow["artifacts"]) {
			artifact := mapValue(raw)
			path, expected := stringValue(artifact["path"]), stringValue(artifact["sha256"])
			actual, err := shaFile(path)
			if err != nil || actual != expected {
				return fmt.Errorf("external runtime digest mismatch: %s", path)
			}
		}
	}
	return nil
}
func dependencyPairsFor(run string, call Call, callMap map[string]Call) ([]dependencyPair, error) {
	pairs := []dependencyPair{}
	for _, id := range call.DependsOn {
		response, err := successfulResponse(run, id)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, dependencyPair{Call: callMap[id], Response: response})
	}
	return pairs, nil
}
func writeText(path, value string) error { return atomicWrite(path, []byte(value), 0o644) }

func (h *Harness) executeCall(run string, call Call, callMap map[string]Call, config Config, freeze Freeze) map[string]any {
	callDir := filepath.Join(run, "calls", call.CallID)
	_ = os.MkdirAll(callDir, 0o755)
	started := utcNow()
	var finalRecord map[string]any
	err := h.withWorktree(run, call.CallID, freeze.SourceCommit, func(worktree string) error {
		kind := stringValue(call.PromptSpec["kind"])
		var prompt string
		var err error
		if kind == "fuser" || kind == "chain" {
			var deps []dependencyPair
			deps, err = dependencyPairsFor(run, call, callMap)
			if err == nil {
				prompt, err = dynamicPrompt(worktree, call, deps)
			}
		} else {
			prompt, err = staticPrompt(worktree, call)
		}
		if err != nil {
			return err
		}
		system, err := readText(worktree, "experiment/prompts/system.txt")
		if err != nil {
			return err
		}
		var schema map[string]any
		if err := readJSON(filepath.Join(worktree, "experiment/schema/findings.schema.json"), &schema); err != nil {
			return err
		}
		packet, err := packetText(worktree, call.Packet)
		if err != nil {
			return err
		}
		schemaCanonical, _ := json.Marshal(schema)
		request := map[string]any{"schema_version": 2, "call_id": call.CallID, "created_at": started, "source_commit": freeze.SourceCommit, "provider": providerMap(call.Provider), "phase": call.Phase, "packet_sha256": shaText(packet), "system_prompt": system, "system_prompt_sha256": shaText(system), "prompt": prompt, "prompt_sha256": shaText(prompt), "schema": schema, "schema_sha256": shaBytes(schemaCanonical), "answer_key_in_context": false, "isolation": freeze.Isolation}
		if err := writeJSON(filepath.Join(callDir, "request.json"), request); err != nil {
			return err
		}
		maxAttempts := config.MaxAttempts
		if maxAttempts < 1 {
			maxAttempts = 1
		}
		var result Invocation
		attemptUsed := 0
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			attemptUsed = attempt
			attemptDir := filepath.Join(callDir, "attempts", fmt.Sprint(attempt))
			if err := os.MkdirAll(attemptDir, 0o755); err != nil {
				return err
			}
			privateBase := filepath.Join(h.Experiment, "private", "sandboxes", filepath.Base(run), call.CallID, fmt.Sprint(attempt))
			_ = os.RemoveAll(privateBase)
			result, err = invoke(privateBase, call.Provider, prompt, system, schema, config.RequestTimeoutSeconds, "findings", nil)
			_ = os.RemoveAll(privateBase)
			if err != nil {
				return err
			}
			if err := writeText(filepath.Join(attemptDir, "stdout.txt"), result.Stdout); err != nil {
				return err
			}
			if err := writeText(filepath.Join(attemptDir, "stderr.txt"), result.Stderr); err != nil {
				return err
			}
			metadata := map[string]any{"attempt": attempt, "command": result.CommandDisplay, "returncode": result.ReturnCode, "latency_seconds": result.LatencySeconds, "parse_error": nil, "usage": result.Usage, "sandbox": result.Sandbox}
			if result.ParseError != "" {
				metadata["parse_error"] = result.ParseError
			}
			if err := writeJSON(filepath.Join(attemptDir, "metadata.json"), metadata); err != nil {
				return err
			}
			if result.ReturnCode == 0 {
				break
			}
			if attempt < maxAttempts {
				backoff := config.RetryBackoffSeconds
				if backoff == 0 {
					backoff = 5
				}
				time.Sleep(time.Duration(backoff * float64(attempt) * float64(time.Second)))
			}
		}
		status := "success"
		if result.ReturnCode != 0 {
			status = "error"
		} else if result.Parsed == nil {
			status = "malformed"
		}
		finalRecord = map[string]any{"schema_version": 2, "call_id": call.CallID, "status": status, "started_at": started, "completed_at": utcNow(), "attempts": attemptUsed, "returncode": result.ReturnCode, "latency_seconds": result.LatencySeconds, "usage": result.Usage, "parse_error": nil, "parsed_response": result.Parsed, "request_sha256": request["prompt_sha256"], "source_commit": freeze.SourceCommit, "worktree": map[string]any{"fresh": true, "detached": true, "removed_after_call": true, "visible_to_model_process": false}, "sandbox": result.Sandbox}
		if result.ParseError != "" {
			finalRecord["parse_error"] = result.ParseError
		}
		return writeJSON(filepath.Join(callDir, "record.json"), finalRecord)
	})
	if err != nil {
		finalRecord = map[string]any{"schema_version": 2, "call_id": call.CallID, "status": "error", "started_at": started, "completed_at": utcNow(), "attempts": 0, "error_type": "GoError", "error": err.Error(), "source_commit": freeze.SourceCommit}
		_ = writeJSON(filepath.Join(callDir, "record.json"), finalRecord)
	}
	return finalRecord
}
func blockedRecord(run string, call Call, failed []string) error {
	return writeJSON(filepath.Join(run, "calls", call.CallID, "record.json"), map[string]any{"schema_version": 2, "call_id": call.CallID, "status": "blocked", "completed_at": utcNow(), "failed_dependencies": failed})
}
func orderedPending(run string, plan Plan) ([]string, error) {
	out := []string{}
	for _, call := range plan.Calls {
		record, err := readCallRecord(run, call.CallID)
		if err != nil {
			return nil, err
		}
		if record == nil || !terminalStatuses[stringValue(record["status"])] {
			out = append(out, call.CallID)
		}
	}
	return out, nil
}

type callResult struct {
	id     string
	record map[string]any
}

func (h *Harness) Run(runArgument string, jobs int) error {
	if err := currentUserCheck(); err != nil {
		return err
	}
	run, err := h.resolveRun(runArgument)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(run, "seal.json")); err == nil {
		return fmt.Errorf("run is sealed and cannot be modified")
	}
	var freeze Freeze
	var config Config
	var plan Plan
	if err := readJSON(filepath.Join(run, "freeze.json"), &freeze); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(run, "config.json"), &config); err != nil {
		return err
	}
	if err := readJSON(filepath.Join(run, "plan.json"), &plan); err != nil {
		return err
	}
	if err := h.verifyManifest(freeze); err != nil {
		return err
	}
	callMap := map[string]Call{}
	for _, call := range plan.Calls {
		callMap[call.CallID] = call
	}
	pending, err := orderedPending(run, plan)
	if err != nil {
		return err
	}
	if jobs <= 0 {
		jobs = config.MaxConcurrency
	}
	if jobs < 1 {
		jobs = 1
	}
	results := make(chan callResult, jobs)
	running := 0
	for len(pending) > 0 || running > 0 {
		progress := false
		for i := 0; i < len(pending) && running < jobs; {
			id := pending[i]
			call := callMap[id]
			ready := true
			failed := []string{}
			for _, dep := range call.DependsOn {
				record, err := readCallRecord(run, dep)
				if err != nil {
					return err
				}
				if record == nil {
					ready = false
					break
				}
				if stringValue(record["status"]) != "success" {
					failed = append(failed, dep)
				}
			}
			if !ready {
				i++
				continue
			}
			pending = append(pending[:i], pending[i+1:]...)
			progress = true
			if len(failed) > 0 {
				if err := blockedRecord(run, call, failed); err != nil {
					return err
				}
				fmt.Printf("%s blocked by %s\n", id, strings.Join(failed, ", "))
				continue
			}
			running++
			go func(c Call) {
				results <- callResult{id: c.CallID, record: h.executeCall(run, c, callMap, config, freeze)}
			}(call)
		}
		if running > 0 {
			result := <-results
			running--
			fmt.Printf("%s %s\n", result.id, stringValue(result.record["status"]))
			progress = true
		}
		if !progress && len(pending) > 0 {
			return fmt.Errorf("plan contains a dependency cycle or missing call")
		}
	}
	statuses := map[string]int{}
	ids := make([]string, 0, len(callMap))
	for id := range callMap {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		record, err := readCallRecord(run, id)
		if err != nil {
			return err
		}
		status := "missing"
		if record != nil {
			status = stringValue(record["status"])
		}
		statuses[status]++
	}
	if err := writeJSON(filepath.Join(run, "run-summary.json"), map[string]any{"completed_at": utcNow(), "statuses": statuses}); err != nil {
		return err
	}
	encoded, _ := json.Marshal(statuses)
	fmt.Println(string(encoded))
	if statuses["missing"] > 0 {
		return fmt.Errorf("run has missing calls")
	}
	return nil
}
