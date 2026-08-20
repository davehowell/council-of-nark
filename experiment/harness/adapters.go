package harness

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type adapterCommand struct{ Args, Display, Executables, RuntimeRoots []string }

func compactJSON(value any) string { data, _ := json.Marshal(value); return string(data) }
func commandFor(provider Provider, prompt, system string, schema map[string]any) (adapterCommand, error) {
	cliSchema := map[string]any{}
	for k, v := range schema {
		if k != "$schema" {
			cliSchema[k] = v
		}
	}
	schemaText := compactJSON(cliSchema)
	effort := provider.Effort
	if effort == "" {
		effort = "low"
	}
	var args, executables, roots []string
	switch provider.Adapter {
	case "claude":
		return adapterCommand{}, fmt.Errorf("direct Claude CLI is disabled in the strict harness: login state is incompatible with an ephemeral HOME; use an explicitly pinned Anthropic model through Pi")
	case "agy":
		return adapterCommand{}, fmt.Errorf("agy is disabled in the strict harness: its OAuth/keychain state is incompatible with an ephemeral HOME; use an explicitly pinned Pi Gemini model")
	case "pi":
		piPath, err := exec.LookPath("pi")
		if err != nil {
			return adapterCommand{}, err
		}
		script, err := filepath.EvalSymlinks(piPath)
		if err != nil {
			return adapterCommand{}, err
		}
		nodePath, err := exec.LookPath("node")
		if err != nil {
			return adapterCommand{}, err
		}
		nodePath, err = filepath.EvalSymlinks(nodePath)
		if err != nil {
			return adapterCommand{}, err
		}
		args = []string{nodePath, script, "--print", "--no-session", "--no-tools", "--no-extensions", "--no-skills", "--no-prompt-templates", "--no-themes", "--no-context-files", "--no-approve", "--model", provider.Model, "--thinking", effort, "--system-prompt", system, "--mode", "json", prompt}
		executables = []string{nodePath}
		roots = []string{filepath.Dir(filepath.Dir(nodePath))}
	case "mock":
		return adapterCommand{}, nil
	default:
		return adapterCommand{}, fmt.Errorf("unknown adapter %s", provider.Adapter)
	}
	display := make([]string, len(args))
	for i, item := range args {
		switch item {
		case prompt:
			display[i] = "<PROMPT>"
		case system:
			display[i] = "<SYSTEM_PROMPT>"
		case schemaText:
			display[i] = "<JSON_SCHEMA>"
		default:
			display[i] = item
		}
	}
	return adapterCommand{Args: args, Display: display, Executables: executables, RuntimeRoots: roots}, nil
}

func decodeJSONPrefix(text string) (any, bool) {
	dec := json.NewDecoder(strings.NewReader(text))
	dec.UseNumber()
	var value any
	if err := dec.Decode(&value); err != nil {
		return nil, false
	}
	return value, true
}

var fencePattern = regexp.MustCompile(`(?is)^\x60\x60\x60(?:json)?\s*|\s*\x60\x60\x60$`)

func rawJSONValues(text string) []any {
	values := []any{}
	trimmed := strings.TrimSpace(text)
	if strings.HasPrefix(trimmed, "```") {
		trimmed = fencePattern.ReplaceAllString(trimmed, "")
	}
	candidates := append([]string{trimmed}, strings.Split(text, "\n")...)
	for _, candidate := range candidates {
		if value, ok := decodeJSONPrefix(candidate); ok {
			values = append(values, value)
		}
	}
	for i, r := range text {
		if r == '{' || r == '[' {
			if value, ok := decodeJSONPrefix(text[i:]); ok {
				values = append(values, value)
			}
		}
	}
	return values
}
func nestedCandidates(value any) []any {
	out := []any{value}
	switch v := value.(type) {
	case map[string]any:
		for _, key := range []string{"result", "structured_output", "output", "response", "message"} {
			if nested, ok := v[key]; ok {
				if text, ok := nested.(string); ok {
					out = append(out, rawJSONValues(text)...)
				} else {
					out = append(out, nestedCandidates(nested)...)
				}
			}
		}
		if content, ok := v["content"].([]any); ok {
			for _, raw := range content {
				block := mapValue(raw)
				if text, ok := block["text"].(string); ok {
					out = append(out, rawJSONValues(text)...)
				}
			}
		}
	case []any:
		for _, item := range v {
			out = append(out, nestedCandidates(item)...)
		}
	}
	return out
}
func extract(text, root string) (map[string]any, string) {
	candidates := []any{}
	for _, value := range rawJSONValues(text) {
		candidates = append(candidates, nestedCandidates(value)...)
	}
	for i := len(candidates) - 1; i >= 0; i-- {
		if value, ok := candidates[i].(map[string]any); ok {
			if _, ok := value[root].([]any); ok {
				return value, ""
			}
		}
	}
	return nil, fmt.Sprintf("no JSON object with array root key %q found", root)
}
func ndjsonEvents(text string) []map[string]any {
	events := []map[string]any{}
	for _, line := range strings.Split(text, "\n") {
		if value, ok := decodeJSONPrefix(line); ok {
			if event, ok := value.(map[string]any); ok {
				events = append(events, event)
			}
		}
	}
	return events
}
func messageText(raw any) []string {
	message := mapValue(raw)
	if stringValue(message["role"]) != "assistant" {
		return nil
	}
	out := []string{}
	for _, rawBlock := range sliceValue(message["content"]) {
		block := mapValue(rawBlock)
		if stringValue(block["type"]) == "text" {
			if text, ok := block["text"].(string); ok {
				out = append(out, text)
			}
		}
	}
	return out
}
func extractPiAssistant(text, root string) (map[string]any, string) {
	completed, deltas := []string{}, []string{}
	for _, value := range ndjsonEvents(text) {
		kind := stringValue(value["type"])
		if kind == "message_end" || kind == "turn_end" {
			completed = append(completed, messageText(value["message"])...)
		}
		if kind == "agent_end" {
			for _, message := range sliceValue(value["messages"]) {
				completed = append(completed, messageText(message)...)
			}
		}
		event := mapValue(value["assistantMessageEvent"])
		if len(event) == 0 {
			event = mapValue(value["event"])
		}
		if stringValue(event["type"]) == "text_end" {
			if content, ok := event["content"].(string); ok {
				completed = append(completed, content)
			}
		} else if stringValue(event["type"]) == "text_delta" {
			delta := stringValue(event["delta"])
			if delta == "<nil>" || delta == "" {
				delta = stringValue(event["text"])
			}
			if delta != "" && delta != "<nil>" {
				deltas = append(deltas, delta)
			}
		}
	}
	if len(deltas) > 0 {
		completed = append(completed, strings.Join(deltas, ""))
	}
	for i := len(completed) - 1; i >= 0; i-- {
		values := rawJSONValues(completed[i])
		for j := len(values) - 1; j >= 0; j-- {
			if value, ok := values[j].(map[string]any); ok {
				if _, ok := value[root].([]any); ok {
					return value, ""
				}
			}
		}
	}
	return nil, fmt.Sprintf("no assistant JSON object with array root key %q found", root)
}
func piProviderError(text string) string {
	for _, value := range ndjsonEvents(text) {
		message := mapValue(value["message"])
		if stringValue(message["role"]) == "assistant" && (stringValue(message["stopReason"]) == "error" || message["errorMessage"] != nil) {
			if s := stringValue(message["errorMessage"]); s != "" && s != "<nil>" {
				return s
			}
			return "Pi provider error"
		}
		if stringValue(value["type"]) == "error" {
			for _, key := range []string{"error", "message"} {
				if s := stringValue(value[key]); s != "" && s != "<nil>" {
					return s
				}
			}
			return "Pi provider error"
		}
	}
	return ""
}
func usageFrom(text string) map[string]any {
	result := map[string]any{}
	for _, raw := range rawJSONValues(text) {
		value, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		for _, key := range []string{"usage", "modelUsage"} {
			if nested, ok := value[key].(map[string]any); ok {
				result[key] = nested
			}
		}
		for _, key := range []string{"total_cost_usd", "cost_usd", "duration_ms", "duration_api_ms"} {
			if value[key] != nil {
				result[key] = value[key]
			}
		}
		if message := mapValue(value["message"]); len(message) > 0 {
			if usage, ok := message["usage"].(map[string]any); ok {
				result["usage"] = usage
			}
		}
	}
	return result
}

func invoke(base string, provider Provider, prompt, system string, schema map[string]any, timeout int, expectedRoot string, expectedIDs map[string]bool) (Invocation, error) {
	command, err := commandFor(provider, prompt, system, schema)
	if err != nil {
		return Invocation{}, err
	}
	started := time.Now()
	if provider.Adapter == "mock" {
		parsed := map[string]any{"findings": []any{}}
		return Invocation{ReturnCode: 0, Stdout: "{\"findings\":[]}", Parsed: parsed, LatencySeconds: time.Since(started).Seconds(), Usage: map[string]any{}, Sandbox: map[string]any{"seatbelt": true, "process": "mock-no-child"}}, nil
	}
	sandbox, err := makeSandbox(base, provider, command.Executables, command.RuntimeRoots)
	if err != nil {
		return Invocation{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	args := append([]string{"-f", sandbox.Profile, "--"}, command.Args...)
	cmd := exec.CommandContext(ctx, "/usr/bin/sandbox-exec", args...)
	cmd.Dir = sandbox.CWD
	cmd.Env = sandbox.Environment
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	runErr := cmd.Run()
	returnCode := 0
	if runErr != nil {
		if exit, ok := runErr.(*exec.ExitError); ok {
			returnCode = exit.ExitCode()
		} else {
			returnCode = 1
		}
	}
	if ctx.Err() == context.DeadlineExceeded {
		returnCode = 124
		stderr.WriteString(fmt.Sprintf("\ntimed out after %d seconds", timeout))
	}
	stdoutText, stderrText := stdout.String(), stderr.String()
	var parsed map[string]any
	var parseError string
	if provider.Adapter == "pi" {
		if providerErr := piProviderError(stdoutText); providerErr != "" {
			if returnCode == 0 {
				returnCode = 75
			}
			stderrText += "\nPi provider error: " + providerErr
			parseError = providerErr
		} else {
			parsed, parseError = extractPiAssistant(stdoutText, expectedRoot)
		}
	} else {
		parsed, parseError = extract(stdoutText, expectedRoot)
	}
	if parsed != nil {
		var validation []string
		if expectedRoot == "findings" {
			validation = findingsErrors(parsed)
		} else {
			validation = judgementErrors(parsed, expectedIDs)
		}
		if len(validation) > 0 {
			parseError = strings.Join(validation, "; ")
			parsed = nil
		}
	}
	profileSHA, _ := shaFile(sandbox.Profile)
	return Invocation{Command: command.Args, CommandDisplay: command.Display, ReturnCode: returnCode, Stdout: stdoutText, Stderr: stderrText, LatencySeconds: time.Since(started).Seconds(), Parsed: parsed, ParseError: parseError, Usage: usageFrom(stdoutText), Sandbox: map[string]any{"seatbelt": true, "profile_sha256": profileSHA, "ephemeral_home": true, "empty_cwd": true, "worktree_visible_to_child": false, "network": "provider CLI transport permitted; model tools disabled"}}, nil
}
