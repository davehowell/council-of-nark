package respondent

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestTrackedConfigLoads(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	config, data, err := LoadConfig(root, "experiment/ecological/config/eco-gortex-unicode-tokenizer-tools.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 || config.TaskID != "eco-gortex-unicode-tokenizer" || len(config.ActiveTools) != 5 {
		t.Fatalf("unexpected tracked config: %#v", config)
	}
}

func TestConfigRequiresExactToolSet(t *testing.T) {
	config := Config{
		SchemaVersion: 1, TaskID: "task", SnapshotControllerCommit: "commit",
		SourceTreeSHA256: strings.Repeat("a", 64), ToolPolicy: "policy", PiExtension: "extension",
		Provider: Provider{Model: "model", Thinking: "off"}, DoctorTimeoutSeconds: 10,
		ActiveTools: []string{"source_list", "source_read", "source_search", "run_focused_test", "source_read"},
	}
	if err := config.Validate(); err == nil {
		t.Fatal("duplicate/missing tool set was accepted")
	}
}

func TestChildProfileExposesOnlyLiteralExtension(t *testing.T) {
	profile := ChildProfile("/private/scratch", []string{"/repo/ecological-tools.ts", "/runtime/pi.js"}, []string{"/runtime/node"}, []string{"/runtime"})
	for _, required := range []string{"(deny default)", "(allow network-outbound)", `(literal "/repo/ecological-tools.ts")`, `(subpath "/private/scratch")`} {
		if !strings.Contains(profile, required) {
			t.Fatalf("profile is missing %s:\n%s", required, profile)
		}
	}
	for _, forbidden := range []string{`(subpath "/repo")`, "/source", "/controller", "/bin/sh"} {
		if strings.Contains(profile, forbidden) {
			t.Fatalf("profile unexpectedly contains %s:\n%s", forbidden, profile)
		}
	}
}

func TestRPCHealthAndStateParsing(t *testing.T) {
	healthFrames := strings.Join([]string{
		`{"type":"extension_ui_request","method":"notify","message":"{\"ecological_mediator_health\":true,\"sequence\":1}"}`,
		`{"id":"health","type":"response","command":"prompt","success":true}`,
	}, "\n") + "\n"
	var events bytes.Buffer
	response, notify, extensionErrors, err := awaitHealth(bufio.NewReader(strings.NewReader(healthFrames)), &events)
	if err != nil || !response || !notify || extensionErrors != 0 {
		t.Fatalf("health parse failed: %v %v %d %v", response, notify, extensionErrors, err)
	}
	commandFrame := `{"id":"thinking","type":"response","command":"set_thinking_level","success":true}` + "\n"
	extensionErrors, err = awaitCommand(bufio.NewReader(strings.NewReader(commandFrame)), &events, "thinking")
	if err != nil || extensionErrors != 0 {
		t.Fatalf("command parse failed: %d %v", extensionErrors, err)
	}
	stateFrame := `{"id":"state","type":"response","command":"get_state","success":true,"data":{"model":{"id":"gemma"},"thinkingLevel":"off"}}` + "\n"
	model, thinking, extensionErrors, err := awaitState(bufio.NewReader(strings.NewReader(stateFrame)), &events)
	if err != nil || model != "gemma" || thinking != "off" || extensionErrors != 0 {
		t.Fatalf("state parse failed: %s %s %d %v", model, thinking, extensionErrors, err)
	}
}

func TestChildBoundaryProbe(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.Skip("Seatbelt is maintained only on macOS")
	}
	runtimeInfo, err := ResolvePiRuntime()
	if err != nil {
		t.Fatal(err)
	}
	extensionRoot := t.TempDir()
	extension := filepath.Join(extensionRoot, "extension.ts")
	if err := os.WriteFile(extension, []byte("export default function () {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	denied := filepath.Join(t.TempDir(), "denied.txt")
	if err := os.WriteFile(denied, []byte("secret\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	sandbox, err := PrepareChildSandbox(t.TempDir(), runtimeInfo, extension)
	if err != nil {
		t.Fatal(err)
	}
	result, err := ProbeChildBoundary(sandbox, runtimeInfo, extension, []string{denied})
	if err != nil {
		t.Fatalf("probe failed: %#v %v", result, err)
	}
}
