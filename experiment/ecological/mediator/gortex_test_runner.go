package mediator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

type GortexTestRunner struct {
	Attempt      string
	Ecological   string
	Controller   string
	ParentRoot   string
	Toolchain    string
	Modules      string
	GoBinary     string
	Package      string
	RunPattern   string
	FailureText  string
	ArtifactRoot string

	mu    sync.Mutex
	calls int
}

type closureFile struct {
	ToolchainPath string `json:"toolchain_path"`
	ToolchainTree string `json:"toolchain_tree_sha256"`
	ModulePath    string `json:"module_path"`
	ModuleTree    string `json:"module_tree_sha256"`
}

type snapshotConfigFile struct {
	TaskID      string `json:"task_id"`
	FocusedTest struct {
		Package               string `json:"package"`
		Run                   string `json:"run"`
		ParentFailureContains string `json:"parent_failure_contains"`
	} `json:"focused_test"`
}

func NewGortexTestRunner(attempt, artifactRoot string) (*GortexTestRunner, error) {
	attempt, err := filepath.Abs(attempt)
	if err != nil {
		return nil, err
	}
	controller := filepath.Join(attempt, "controller")
	var closure closureFile
	if err := readJSONFile(filepath.Join(controller, "closure.json"), &closure); err != nil {
		return nil, fmt.Errorf("load snapshot closure: %w", err)
	}
	var config snapshotConfigFile
	if err := readJSONFile(filepath.Join(controller, "config.json"), &config); err != nil {
		return nil, fmt.Errorf("load snapshot config: %w", err)
	}
	if config.TaskID != "eco-gortex-unicode-tokenizer" {
		return nil, fmt.Errorf("Gortex runner does not support task %q", config.TaskID)
	}
	ecological := filepath.Dir(filepath.Dir(attempt))
	toolchain := filepath.Join(ecological, filepath.FromSlash(closure.ToolchainPath))
	modules := filepath.Join(ecological, filepath.FromSlash(closure.ModulePath))
	for path, expected := range map[string]string{toolchain: closure.ToolchainTree, modules: closure.ModuleTree} {
		manifest, err := snapshot.BuildClosureManifest(path)
		if err != nil {
			return nil, fmt.Errorf("verify frozen closure %s: %w", path, err)
		}
		if manifest.TreeSHA256 != expected {
			return nil, fmt.Errorf("frozen closure digest mismatch: %s", path)
		}
	}
	goBinary := filepath.Join(toolchain, "root", "bin", "go")
	if info, err := os.Stat(goBinary); err != nil || info.Mode()&0o111 == 0 {
		return nil, fmt.Errorf("frozen Go executable is unavailable")
	}
	parentRoot := filepath.Join(controller, "validation", "parent")
	if info, err := os.Stat(parentRoot); err != nil || !info.IsDir() {
		return nil, fmt.Errorf("validated parent test root is unavailable")
	}
	if artifactRoot == "" {
		return nil, fmt.Errorf("test artifact root is required")
	}
	artifactRoot, err = filepath.Abs(artifactRoot)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(artifactRoot, 0o700); err != nil {
		return nil, err
	}
	return &GortexTestRunner{
		Attempt: attempt, Ecological: ecological, Controller: controller, ParentRoot: parentRoot,
		Toolchain: toolchain, Modules: modules, GoBinary: goBinary,
		Package: config.FocusedTest.Package, RunPattern: config.FocusedTest.Run,
		FailureText: config.FocusedTest.ParentFailureContains, ArtifactRoot: artifactRoot,
	}, nil
}

func (r *GortexTestRunner) Run(target string) (string, map[string]any, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if target != "unicode-tokenizer-regression" {
		return "", nil, fmt.Errorf("unsupported Gortex test target %q", target)
	}
	r.calls++
	callRoot := filepath.Join(r.ArtifactRoot, fmt.Sprintf("%02d", r.calls))
	scratch := filepath.Join(callRoot, "scratch")
	if err := os.MkdirAll(scratch, 0o700); err != nil {
		return "", nil, err
	}
	profileText := gortexTestProfile(r.ParentRoot, filepath.Join(r.Toolchain, "root"), filepath.Join(r.Modules, "gomod"), r.GoBinary, scratch)
	profilePath := filepath.Join(callRoot, "profile.sb")
	if err := os.WriteFile(profilePath, []byte(profileText), 0o600); err != nil {
		return "", nil, err
	}
	env := gortexTestEnvironment(r.Toolchain, r.Modules, scratch)
	probe := runNetworkProbe(profilePath, scratch, env)
	if err := writeJSONFile(filepath.Join(callRoot, "network-probe.json"), probe.log); err != nil {
		return "", nil, err
	}
	if probe.err != nil {
		return "", nil, probe.err
	}
	executableProbe := runSandboxCommand(profilePath, scratch, env, "/bin/sh", "-c", "true")
	if err := writeJSONFile(filepath.Join(callRoot, "executable-probe.json"), map[string]any{
		"command":   []string{"sandbox-exec", "<test-profile>", "sh", "<probe>"},
		"exit_code": executableProbe.exitCode, "stdout": executableProbe.stdout, "stderr": executableProbe.stderr,
		"unlisted_executable_denied": executableProbe.exitCode != 0,
	}); err != nil {
		return "", nil, err
	}
	if executableProbe.exitCode == 0 {
		return "", nil, fmt.Errorf("focused-test sandbox executed an unlisted shell")
	}

	args := []string{"test", "-count=1", r.Package, "-run", r.RunPattern}
	started := time.Now()
	result := runSandboxCommand(profilePath, r.ParentRoot, env, r.GoBinary, args...)
	duration := time.Since(started)
	log := map[string]any{
		"command": append([]string{"<frozen-go>"}, args...), "exit_code": result.exitCode,
		"duration_seconds": duration.Seconds(), "stdout": result.stdout, "stderr": result.stderr,
	}
	if err := writeJSONFile(filepath.Join(callRoot, "test.json"), log); err != nil {
		return "", nil, err
	}
	rawOutput := strings.TrimSpace(result.stdout + "\n" + result.stderr)
	if result.exitCode == 0 || !strings.Contains(rawOutput, r.FailureText) {
		return "", nil, fmt.Errorf("focused parent test did not reproduce the frozen failure (exit %d)", result.exitCode)
	}
	safeOutput := r.sanitizeOutput(rawOutput, scratch)
	profileDigest, err := snapshot.FileSHA256(profilePath)
	if err != nil {
		return "", nil, err
	}
	return "The focused regression reproduced the frozen parent failure:\n\n" + safeOutput, map[string]any{
		"target": target, "exit_code": result.exitCode, "expected_failure_observed": true,
		"network_probe_denied": true, "unlisted_executable_denied": true, "profile_sha256": profileDigest,
		"duration_seconds": duration.Seconds(),
	}, nil
}

func (r *GortexTestRunner) sanitizeOutput(output, scratch string) string {
	replacements := [][2]string{
		{r.ParentRoot, "<TEST_ROOT>"},
		{filepath.Join(r.Toolchain, "root"), "<FROZEN_GOROOT>"},
		{filepath.Join(r.Modules, "gomod"), "<FROZEN_GOMODCACHE>"},
		{scratch, "<TEST_SCRATCH>"},
		{r.Ecological, "<ECOLOGICAL_ROOT>"},
	}
	for _, replacement := range replacements {
		output = strings.ReplaceAll(output, replacement[0], replacement[1])
	}
	return output
}

type commandOutput struct {
	exitCode int
	stdout   string
	stderr   string
}

func runSandboxCommand(profile, cwd string, env []string, executable string, args ...string) commandOutput {
	sandboxArgs := []string{"-f", profile, "--", executable}
	sandboxArgs = append(sandboxArgs, args...)
	command := exec.Command("/usr/bin/sandbox-exec", sandboxArgs...)
	command.Dir = cwd
	command.Env = env
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	exit := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exit = exitError.ExitCode()
		} else {
			exit = -1
			stderr.WriteString("\nlauncher error: " + err.Error())
		}
	}
	return commandOutput{exitCode: exit, stdout: stdout.String(), stderr: stderr.String()}
}

type networkProbeResult struct {
	log map[string]any
	err error
}

func runNetworkProbe(profile, cwd string, env []string) networkProbeResult {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return networkProbeResult{err: err}
	}
	defer listener.Close()
	if tcp, ok := listener.(*net.TCPListener); ok {
		_ = tcp.SetDeadline(time.Now().Add(3 * time.Second))
	}
	port := listener.Addr().(*net.TCPAddr).Port
	accepted := make(chan bool, 1)
	go func() {
		connection, acceptErr := listener.Accept()
		if acceptErr == nil {
			_ = connection.Close()
			accepted <- true
			return
		}
		accepted <- false
	}()
	started := time.Now()
	result := runSandboxCommand(profile, cwd, env, "/usr/bin/nc", "-z", "-w", "1", "127.0.0.1", strconv.Itoa(port))
	_ = listener.Close()
	connected := <-accepted
	log := map[string]any{
		"command":   []string{"sandbox-exec", "<test-profile>", "nc", "<controller-listener>"},
		"exit_code": result.exitCode, "duration_seconds": time.Since(started).Seconds(),
		"stdout": result.stdout, "stderr": result.stderr, "controller_connected": connected,
	}
	if result.exitCode == 0 || connected {
		return networkProbeResult{log: log, err: fmt.Errorf("focused-test sandbox reached a controller listener")}
	}
	return networkProbeResult{log: log}
}

func gortexTestProfile(parentRoot, goroot, gomod, goBinary, scratch string) string {
	readRoots := []string{
		"/System", "/usr/lib", "/usr/share", "/Library/Apple", "/private/etc",
		"/private/var/db/timezone", "/dev", parentRoot, goroot, gomod, scratch,
	}
	sort.Strings(readRoots)
	var builder strings.Builder
	builder.WriteString("(version 1)\n(deny default)\n")
	builder.WriteString("(allow file-read-metadata)\n(allow sysctl-read)\n(allow mach-lookup)\n(allow system-socket)\n(allow process-fork)\n")
	builder.WriteString("(allow file-read*\n  (literal \"/\")\n")
	for _, path := range readRoots {
		builder.WriteString("  (subpath " + strconv.Quote(path) + ")\n")
	}
	for _, path := range []string{goBinary, "/usr/bin/nc"} {
		builder.WriteString("  (literal " + strconv.Quote(path) + ")\n")
	}
	builder.WriteString(")\n")
	builder.WriteString("(allow file-write* (subpath " + strconv.Quote(scratch) + ") (literal \"/dev/null\"))\n")
	builder.WriteString("(allow process-exec\n")
	for _, path := range []string{goBinary, "/usr/bin/nc"} {
		builder.WriteString("  (literal " + strconv.Quote(path) + ")\n")
	}
	for _, path := range []string{filepath.Join(goroot, "pkg", "tool"), scratch} {
		builder.WriteString("  (subpath " + strconv.Quote(path) + ")\n")
	}
	builder.WriteString(")\n")
	return builder.String()
}

func gortexTestEnvironment(toolchain, modules, scratch string) []string {
	values := map[string]string{
		"HOME": filepath.Join(scratch, "home"), "TMPDIR": filepath.Join(scratch, "tmp"),
		"GOTMPDIR": filepath.Join(scratch, "gotmp"), "GOCACHE": filepath.Join(scratch, "gocache"),
		"GOMODCACHE": filepath.Join(modules, "gomod"), "GOROOT": filepath.Join(toolchain, "root"),
		"GOENV": "off", "GOWORK": "off", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0",
		"GOFLAGS": "-buildvcs=false", "GOSUMDB": "off", "GOPROXY": "off",
		"PATH": filepath.Join(toolchain, "root", "bin") + ":/usr/bin:/bin", "LANG": "en_US.UTF-8",
	}
	for _, path := range []string{values["HOME"], values["TMPDIR"], values["GOTMPDIR"], values["GOCACHE"]} {
		_ = os.MkdirAll(path, 0o700)
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	env := make([]string, 0, len(keys))
	for _, key := range keys {
		env = append(env, key+"="+values[key])
	}
	return env
}

func readJSONFile(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func writeJSONFile(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o600)
}
