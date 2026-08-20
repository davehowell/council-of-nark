package harness

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func testHarness(t *testing.T) *Harness {
	t.Helper()
	h, err := New()
	if err != nil {
		t.Fatal(err)
	}
	return h
}
func configAt(t *testing.T, name string) Config {
	t.Helper()
	h := testHarness(t)
	var config Config
	if err := readJSON(filepath.Join(h.Root, "experiment/config", name), &config); err != nil {
		t.Fatal(err)
	}
	return config
}
func TestPlanCounts(t *testing.T) {
	cases := []struct {
		name        string
		calls, sets int
	}{{"stage-a-smoke.json", 81, 27}, {"topology-smoke.json", 144, 108}, {"provider-pair-smoke.json", 18, 18}, {"persona-pair-gemma-repeated.json", 60, 60}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := BuildPlan(configAt(t, tc.name))
			if err != nil {
				t.Fatal(err)
			}
			if len(plan.Calls) != tc.calls || len(plan.OutputSets) != tc.sets {
				t.Fatalf("got %d calls/%d sets", len(plan.Calls), len(plan.OutputSets))
			}
		})
	}
}
func TestPlanOrderDeterministic(t *testing.T) {
	config := configAt(t, "stage-a-smoke.json")
	a, _ := BuildPlan(config)
	b, _ := BuildPlan(config)
	for i := range a.Calls {
		if a.Calls[i].CallID != b.Calls[i].CallID {
			t.Fatal("call order changed")
		}
	}
}
func TestStaticPromptsExcludeKeys(t *testing.T) {
	h := testHarness(t)
	plan, err := BuildPlan(configAt(t, "stage-a-smoke.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, call := range plan.Calls {
		kind := stringValue(call.PromptSpec["kind"])
		if kind == "fuser" || kind == "chain" {
			continue
		}
		prompt, err := staticPrompt(h.Root, call)
		if err != nil {
			t.Fatal(err)
		}
		answer, _ := answerText(h.Root, call.Packet)
		if err := assertAnswerKeyAbsent(prompt, answer); err != nil {
			t.Fatalf("%s: %v", call.CallID, err)
		}
	}
}
func TestExtractRejectsEmbeddedSchemaRoot(t *testing.T) {
	text := `{"json_schema":{"properties":{"judgements":{"type":"array"}}},"structured_output":{"judgements":[]}}`
	value, errText := extract(text, "judgements")
	if errText != "" || value == nil {
		t.Fatalf("extract failed: %s", errText)
	}
	if _, ok := value["judgements"].([]any); !ok {
		t.Fatal("did not select output array")
	}
}
func TestPiAssistantParserIgnoresEchoedPrompt(t *testing.T) {
	text := "{\"type\":\"message_end\",\"message\":{\"role\":\"user\",\"content\":[{\"type\":\"text\",\"text\":\"{\\\"findings\\\":[{\\\"bad\\\":true}]}\"}]}}\n" + "{\"type\":\"message_end\",\"message\":{\"role\":\"assistant\",\"content\":[{\"type\":\"text\",\"text\":\"{\\\"findings\\\":[]}\"}]}}"
	value, errText := extractPiAssistant(text, "findings")
	if errText != "" || len(sliceValue(value["findings"])) != 0 {
		t.Fatalf("unexpected parse: %v %s", value, errText)
	}
}
func TestAgyModelEncodesEffort(t *testing.T) {
	command, err := commandFor(Provider{Adapter: "agy", Model: "gemini-3.5-flash-low", Effort: "low"}, "prompt", "system", map[string]any{"type": "object"})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(command.Args, " ")
	if strings.Contains(joined, "--effort") {
		t.Fatal("agy effort must not be passed separately from its model ID")
	}
	if !strings.Contains(joined, "--model gemini-3.5-flash-low") {
		t.Fatal("agy model was not pinned")
	}
}

func TestGemmaCommandPinsIsolation(t *testing.T) {
	command, err := commandFor(Provider{Adapter: "pi", Model: "gemma-4-31b-it", Effort: "off"}, "prompt", "system", map[string]any{"$schema": "ignored", "type": "object"})
	if err != nil {
		t.Fatal(err)
	}
	joined := strings.Join(command.Args, " ")
	for _, required := range []string{"--no-session", "--no-tools", "--no-extensions", "--no-context-files", "--model gemma-4-31b-it", "--thinking off"} {
		if !strings.Contains(joined, required) {
			t.Errorf("missing %s", required)
		}
	}
	if strings.Contains(joined, "$schema") {
		t.Fatal("CLI schema retained unsupported annotation")
	}
}
func TestBootstrapDeterministic(t *testing.T) {
	a := bootstrapMeanCI([]float64{-.1, 0, .1, .2}, "fixed", 1000)
	b := bootstrapMeanCI([]float64{-.1, 0, .1, .2}, "fixed", 1000)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("bootstrap changed")
	}
}
func TestSeatbeltProbe(t *testing.T) {
	if err := testHarness(t).sandboxProbe(false); err != nil {
		t.Fatal(err)
	}
}
func TestProfileDoesNotAllowRepository(t *testing.T) {
	h := testHarness(t)
	profile := profileText("/private/tmp/isolated", []string{"/bin/echo"}, nil, false)
	if strings.Contains(profile, h.Root) || strings.Contains(profile, filepath.Dir(h.Root)) {
		t.Fatal("profile exposes repository")
	}
	if !strings.Contains(profile, "(deny default)") {
		t.Fatal("profile is not deny-by-default")
	}
}
