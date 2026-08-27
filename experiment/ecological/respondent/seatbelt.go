package respondent

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type PiRuntime struct {
	Node         string
	Script       string
	RuntimeRoots []string
}

type ChildSandbox struct {
	Root        string
	Home        string
	Temp        string
	CWD         string
	Profile     string
	Environment []string
}

var authEnvironment = map[string]bool{
	"ANTHROPIC_API_KEY": true, "ANTHROPIC_OAUTH_TOKEN": true, "OPENAI_API_KEY": true,
	"GEMINI_API_KEY": true, "GOOGLE_API_KEY": true, "GOOGLE_APPLICATION_CREDENTIALS": true,
	"GOOGLE_CLOUD_PROJECT": true, "GOOGLE_CLOUD_LOCATION": true, "AWS_PROFILE": true, "AWS_REGION": true,
}
var baseEnvironment = map[string]bool{"LANG": true, "LC_ALL": true, "SHELL": true, "TERM": true, "USER": true}

func ResolvePiRuntime() (PiRuntime, error) {
	piPath, err := exec.LookPath("pi")
	if err != nil {
		return PiRuntime{}, err
	}
	script, err := filepath.EvalSymlinks(piPath)
	if err != nil {
		return PiRuntime{}, err
	}
	nodePath, err := exec.LookPath("node")
	if err != nil {
		return PiRuntime{}, err
	}
	nodePath, err = filepath.EvalSymlinks(nodePath)
	if err != nil {
		return PiRuntime{}, err
	}
	return PiRuntime{Node: nodePath, Script: script, RuntimeRoots: []string{filepath.Dir(filepath.Dir(nodePath))}}, nil
}

func PrepareChildSandbox(base string, runtime PiRuntime, extension string) (ChildSandbox, error) {
	base, err := filepath.Abs(base)
	if err != nil {
		return ChildSandbox{}, err
	}
	if canonical, err := filepath.EvalSymlinks(base); err == nil {
		base = canonical
	}
	extension, err = filepath.EvalSymlinks(extension)
	if err != nil {
		return ChildSandbox{}, err
	}
	root := filepath.Join(base, "pi-sandbox")
	home := filepath.Join(root, "home")
	temp := filepath.Join(root, "tmp")
	cwd := filepath.Join(root, "cwd")
	for _, path := range []string{home, temp, cwd} {
		if err := os.MkdirAll(path, 0o700); err != nil {
			return ChildSandbox{}, err
		}
	}
	environment, err := prepareEphemeralHome(home)
	if err != nil {
		return ChildSandbox{}, err
	}
	environment["TMPDIR"] = temp
	environment["XDG_RUNTIME_DIR"] = temp
	environment["BUN_TMPDIR"] = temp
	environment["CLAUDE_CODE_TMPDIR"] = temp
	environment["PATH"] = "/usr/bin:/bin"
	profile := filepath.Join(root, "profile.sb")
	text := ChildProfile(root, []string{extension, runtime.Script}, []string{runtime.Node}, runtime.RuntimeRoots)
	if err := os.WriteFile(profile, []byte(text), 0o600); err != nil {
		return ChildSandbox{}, err
	}
	return ChildSandbox{Root: root, Home: home, Temp: temp, CWD: cwd, Profile: profile, Environment: sortedEnvironment(environment)}, nil
}

func ChildProfile(scratch string, readFiles, executables, runtimeRoots []string) string {
	readRoots := []string{
		"/System", "/usr/lib", "/usr/share", "/Library/Apple", "/private/etc",
		"/private/var/db/timezone", "/dev", scratch,
	}
	readRoots = append(readRoots, runtimeRoots...)
	readRoots = unique(readRoots)
	readFiles = unique(readFiles)
	executables = unique(executables)
	var builder strings.Builder
	builder.WriteString("(version 1)\n(deny default)\n")
	builder.WriteString("(allow file-read-metadata)\n(allow sysctl-read)\n(allow mach-lookup)\n(allow network-outbound)\n(allow system-socket)\n(allow process*)\n")
	builder.WriteString("(allow file-read*\n  (literal \"/\")\n")
	for _, path := range readRoots {
		builder.WriteString("  (subpath " + strconv.Quote(path) + ")\n")
	}
	for _, path := range readFiles {
		builder.WriteString("  (literal " + strconv.Quote(path) + ")\n")
	}
	for _, path := range executables {
		builder.WriteString("  (literal " + strconv.Quote(path) + ")\n")
	}
	builder.WriteString(")\n")
	builder.WriteString("(allow file-write* (subpath " + strconv.Quote(scratch) + ") (literal \"/dev/null\"))\n")
	builder.WriteString("(allow process-exec\n")
	for _, path := range executables {
		builder.WriteString("  (literal " + strconv.Quote(path) + ")\n")
	}
	builder.WriteString(")\n")
	return builder.String()
}

func ProbeChildBoundary(sandbox ChildSandbox, runtime PiRuntime, extension string, denied []string) (map[string]any, error) {
	var err error
	extension, err = filepath.EvalSymlinks(extension)
	if err != nil {
		return nil, err
	}
	for index, path := range denied {
		if canonical, resolveErr := filepath.EvalSymlinks(path); resolveErr == nil {
			denied[index] = canonical
		}
	}
	probeProfilePath := filepath.Join(sandbox.Root, "probe-profile.sb")
	profileText := ChildProfile(sandbox.Root, []string{extension, runtime.Script}, []string{runtime.Node, "/bin/sh"}, runtime.RuntimeRoots)
	if err := os.WriteFile(probeProfilePath, []byte(profileText), 0o600); err != nil {
		return nil, err
	}
	values := environmentMap(sandbox.Environment)
	values["EXTENSION"] = extension
	for index, path := range denied {
		values[fmt.Sprintf("DENIED_%d", index)] = path
	}
	var script strings.Builder
	script.WriteString("set -eu\nIFS= read -r first < \"$EXTENSION\"\ntest -n \"$first\"\nprintf allowed > \"$TMPDIR/probe-write\"\n")
	for index := range denied {
		script.WriteString(fmt.Sprintf("if IFS= read -r leaked < \"$DENIED_%d\"; then exit %d; fi\n", index, 90+index))
	}
	command := exec.Command("/usr/bin/sandbox-exec", "-f", probeProfilePath, "--", "/bin/sh", "-c", script.String())
	command.Dir = sandbox.CWD
	command.Env = sortedEnvironment(values)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	exit := 0
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exit = exitError.ExitCode()
		} else {
			exit = -1
		}
	}
	result := map[string]any{
		"command":   []string{"sandbox-exec", "<probe-profile>", "sh", "<probe>"},
		"exit_code": exit, "stdout": stdout.String(), "stderr": stderr.String(),
		"extension_read": exit == 0, "scratch_write": exit == 0, "denied_reads": exit == 0,
	}
	if err != nil {
		return result, fmt.Errorf("Pi child boundary probe failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	if data, err := os.ReadFile(filepath.Join(sandbox.Temp, "probe-write")); err != nil || string(data) != "allowed" {
		return result, fmt.Errorf("Pi child boundary did not permit scratch write")
	}
	return result, nil
}

func prepareEphemeralHome(home string) (map[string]string, error) {
	original, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	for _, relative := range []string{".pi/agent/auth.json", ".pi/agent/models-store.json"} {
		if err := copyIfPresent(filepath.Join(original, relative), filepath.Join(home, relative)); err != nil {
			return nil, err
		}
	}
	environment := map[string]string{}
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 && (baseEnvironment[parts[0]] || authEnvironment[parts[0]]) {
			environment[parts[0]] = parts[1]
		}
	}
	if source := environment["GOOGLE_APPLICATION_CREDENTIALS"]; source != "" {
		destination := filepath.Join(home, "credentials", "google-application-credentials.json")
		if err := copyIfPresent(source, destination); err != nil {
			return nil, err
		}
		environment["GOOGLE_APPLICATION_CREDENTIALS"] = destination
	}
	environment["HOME"] = home
	environment["XDG_CONFIG_HOME"] = filepath.Join(home, ".config")
	environment["XDG_CACHE_HOME"] = filepath.Join(home, ".cache")
	environment["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = "1"
	environment["DISABLE_AUTOUPDATER"] = "1"
	environment["PI_SKIP_VERSION_CHECK"] = "1"
	environment["PI_TELEMETRY"] = "0"
	environment["NO_COLOR"] = "1"
	return environment, nil
}

func copyIfPresent(source, destination string) error {
	info, err := os.Stat(source)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("credential source is not a regular file: %s", source)
	}
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}
	return os.WriteFile(destination, data, 0o600)
}

func unique(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

func sortedEnvironment(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, key+"="+values[key])
	}
	return result
}

func environmentMap(values []string) map[string]string {
	result := map[string]string{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
