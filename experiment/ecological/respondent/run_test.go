package respondent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureReview() EcologicalReview {
	return EcologicalReview{
		ReviewSummary: "A supported diagnosis.",
		Findings: []EcologicalFinding{{
			Title: "Boundary error", Locations: []SourceLocation{{Path: "pkg/token.go", LineStart: 10, LineEnd: 12}},
			Claim: "A boundary check uses the wrong unit.", Mechanism: "The index is measured in bytes but applied to decoded values.",
			Consequence: "A valid non-ASCII input can panic.", Correction: "Use a boundary-safe lookahead.",
			AllocationAndScope: "Keep the scan incremental and preserve ASCII behavior.",
			RegressionTests:    []RegressionCase{{Case: "non-ASCII boundary", Expected: "no panic", WhyDiscriminating: "fails at the unsafe boundary"}},
			Evidence:           []string{"pkg/token.go:10-12"}, Confidence: "high",
		}},
		Uncertainties: []string{"Other callers were not exercised."},
	}
}

func TestTrackedRunConfigLoads(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	config, data, err := LoadRunConfig(root, "experiment/ecological/config/eco-gortex-unicode-tokenizer-respondent-mock.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || config.Adapter != "mock" || config.Budgets.MaxTotalTokens != 1000 {
		t.Fatalf("unexpected run config: %#v", config)
	}
}

func TestAnalyzeEventsRequiresOneFinalSubmission(t *testing.T) {
	review := fixtureReview()
	usage := map[string]any{"input": 100, "output": 50, "cacheRead": 20, "cacheWrite": 10, "totalTokens": 180, "cost": map[string]any{"total": 0.25}}
	events := []any{
		map[string]any{"type": "session", "version": 3},
		map[string]any{"type": "agent_start"},
		map[string]any{"type": "message_end", "message": map[string]any{"role": "assistant", "usage": usage}},
		map[string]any{"type": "tool_execution_end", "toolName": "submit_ecological_review", "result": map[string]any{"details": review}, "isError": false},
		map[string]any{"type": "agent_end"},
		map[string]any{"type": "agent_settled"},
	}
	path := filepath.Join(t.TempDir(), "events.jsonl")
	if err := writeJSONLines(path, events); err != nil {
		t.Fatal(err)
	}
	summary, err := AnalyzeEvents(path, 32768)
	if err != nil {
		t.Fatal(err)
	}
	if summary.SubmissionCount != 1 || !summary.SubmissionWasFinalTool || summary.Usage.Total != 180 || summary.Usage.Cost.Total != 0.25 {
		t.Fatalf("unexpected event summary: %#v", summary)
	}

	events = append(events[:4], events[3:]...)
	duplicatePath := filepath.Join(t.TempDir(), "events.jsonl")
	if err := writeJSONLines(duplicatePath, events); err != nil {
		t.Fatal(err)
	}
	if _, err := AnalyzeEvents(duplicatePath, 32768); err == nil || !strings.Contains(err.Error(), "exactly one") {
		t.Fatalf("duplicate submission was accepted: %v", err)
	}
}

func TestAnalyzeAuditRequiresSequenceAndPayload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.jsonl")
	records := []any{
		map[string]any{"schema_version": 1, "sequence": 1, "kind": "session_start"},
		map[string]any{"schema_version": 1, "sequence": 2, "kind": "provider_request", "payload": map[string]any{"model": "mock"}},
		map[string]any{"schema_version": 1, "sequence": 3, "kind": "provider_response"},
		map[string]any{"schema_version": 1, "sequence": 4, "kind": "assistant_usage"},
		map[string]any{"schema_version": 1, "sequence": 5, "kind": "final_submission"},
	}
	if err := writeJSONLines(path, records); err != nil {
		t.Fatal(err)
	}
	summary, err := AnalyzeAudit(path)
	if err != nil || summary.ProviderRequests != 1 || summary.FinalSubmissions != 1 {
		t.Fatalf("valid audit rejected: %#v %v", summary, err)
	}
	records[1] = map[string]any{"schema_version": 1, "sequence": 3, "kind": "provider_request"}
	badPath := filepath.Join(t.TempDir(), "audit.jsonl")
	if err := writeJSONLines(badPath, records); err != nil {
		t.Fatal(err)
	}
	if _, err := AnalyzeAudit(badPath); err == nil {
		t.Fatal("broken audit sequence was accepted")
	}
}

func TestReviewRejectsEscapingLocation(t *testing.T) {
	review := fixtureReview()
	for _, invalid := range []string{"../evidence/key.md", "pkg/../evidence/key.md", "pkg/.git/config", "/tmp/file"} {
		review.Findings[0].Locations[0].Path = invalid
		data, _ := json.Marshal(review)
		if _, err := DecodeReview(data, 32768); err == nil || !strings.Contains(err.Error(), "source snapshot") {
			t.Fatalf("escaping location %q was accepted: %v", invalid, err)
		}
	}
}

func TestRemoveEphemeralRunState(t *testing.T) {
	root := t.TempDir()
	for _, relative := range []string{"pi-sandbox/home/auth.json", "pi-sandbox/tmp/cache", "pi-sandbox/cwd/file", "focused-tests/01/scratch/gocache/item"} {
		full := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(full), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("temporary"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := RemoveEphemeralRunState(root); err != nil {
		t.Fatal(err)
	}
	for _, relative := range []string{"pi-sandbox/home", "pi-sandbox/tmp", "pi-sandbox/cwd", "focused-tests/01/scratch"} {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(relative))); !os.IsNotExist(err) {
			t.Fatalf("ephemeral path remains: %s", relative)
		}
	}
}

func TestSealCoversEveryArtifactAndDetectsMutation(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "status.json"), []byte("{}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "raw"), 0o700); err != nil {
		t.Fatal(err)
	}
	artifact := filepath.Join(root, "raw", "events.jsonl")
	if err := os.WriteFile(artifact, []byte("{\"type\":\"agent_start\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	seal, err := SealAttempt(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(seal.Files) != 2 {
		t.Fatalf("got %d sealed files", len(seal.Files))
	}
	if _, err := VerifyAttemptSeal(root); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(artifact, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(artifact, []byte("tampered\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := VerifyAttemptSeal(root); err == nil {
		t.Fatal("tampered sealed artifact was accepted")
	}
}
