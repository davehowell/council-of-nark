package harness

import (
	"os"
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
	}{{"stage-a-smoke.json", 81, 27}, {"topology-smoke.json", 144, 108}, {"provider-pair-smoke.json", 18, 18}, {"persona-pair-gemma-repeated.json", 60, 60}, {"persona-factorial-gemma.json", 480, 480}}
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
func TestPersonaFactorialBalancesEveryRole(t *testing.T) {
	plan, err := BuildPlan(configAt(t, "persona-factorial-gemma.json"))
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]map[string]int{}
	for _, set := range plan.OutputSets {
		role, wrapper := stringValue(set.Metadata["role"]), stringValue(set.Metadata["wrapper"])
		if counts[role] == nil {
			counts[role] = map[string]int{}
		}
		counts[role][wrapper]++
	}
	if len(counts) != 8 {
		t.Fatalf("got %d roles", len(counts))
	}
	for role, wrappers := range counts {
		if wrappers["functional"] != 30 || wrappers["fictional"] != 30 {
			t.Fatalf("%s: %#v", role, wrappers)
		}
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
func TestSharedLoginAdaptersFailClosed(t *testing.T) {
	for _, adapter := range []string{"agy", "claude"} {
		_, err := commandFor(Provider{Adapter: adapter, Model: "explicit", Effort: "low"}, "prompt", "system", map[string]any{"type": "object"})
		if err == nil || !strings.Contains(err.Error(), "disabled") {
			t.Fatalf("%s: %v", adapter, err)
		}
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
func TestPiEphemeralHomeExcludesSettingsAndSkills(t *testing.T) {
	home := t.TempDir()
	if _, err := prepareEphemeralHome(home, Provider{Adapter: "pi"}); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{".pi/agent/settings.json", ".pi/settings.json", ".pi/agent/skills"} {
		if _, err := os.Stat(filepath.Join(home, rel)); !os.IsNotExist(err) {
			t.Fatalf("unexpected Pi state copied: %s", rel)
		}
	}
}

func TestClaudeTempVariablesPointAtEphemeralScratch(t *testing.T) {
	home := t.TempDir()
	env, err := prepareEphemeralHome(home, Provider{Adapter: "claude"})
	if err != nil {
		t.Fatal(err)
	}
	// makeSandbox sets all four to its per-attempt temp directory; the home preparation
	// must not import any of them from the controller environment.
	for _, key := range []string{"TMPDIR", "XDG_RUNTIME_DIR", "BUN_TMPDIR", "CLAUDE_CODE_TMPDIR"} {
		if _, exists := env[key]; exists {
			t.Fatalf("%s leaked from controller", key)
		}
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
