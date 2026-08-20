package harness

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

type sandboxContext struct {
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

func seatbeltQuote(value string) string { return fmt.Sprintf("%q", value) }
func uniquePaths(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range values {
		if v == "" || seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func profileText(scratch string, executables, runtimeRoots []string, localListener bool) string {
	readSubpaths := []string{"/System", "/usr/lib", "/usr/share", "/Library/Apple", "/private/etc", "/private/var/db/timezone", "/dev", scratch}
	readSubpaths = append(readSubpaths, runtimeRoots...)
	readSubpaths = uniquePaths(readSubpaths)
	var b strings.Builder
	b.WriteString("(version 1)\n(deny default)\n")
	b.WriteString("(allow file-read-metadata)\n(allow sysctl-read)\n(allow mach-lookup)\n(allow network-outbound)\n(allow system-socket)\n(allow process*)\n")
	if localListener {
		b.WriteString("(allow network-bind (local ip))\n(allow network-inbound (local ip))\n")
	}
	b.WriteString("(allow file-read*\n  (literal \"/\")\n")
	for _, p := range readSubpaths {
		b.WriteString("  (subpath " + seatbeltQuote(p) + ")\n")
	}
	for _, p := range uniquePaths(executables) {
		b.WriteString("  (literal " + seatbeltQuote(p) + ")\n")
	}
	b.WriteString(")\n")
	b.WriteString("(allow file-write* (subpath " + seatbeltQuote(scratch) + ") (literal \"/dev/null\"))\n")
	b.WriteString("(allow process-exec\n")
	for _, p := range uniquePaths(executables) {
		b.WriteString("  (literal " + seatbeltQuote(p) + ")\n")
	}
	b.WriteString(")\n")
	return b.String()
}

func copyIfPresent(source, destination string) error {
	if _, err := os.Stat(source); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return copyFile(source, destination, 0o600)
}
func prepareEphemeralHome(home, adapter string) (map[string]string, error) {
	if err := os.MkdirAll(home, 0o700); err != nil {
		return nil, err
	}
	original, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if adapter == "pi" {
		for _, rel := range []string{".pi/agent/auth.json", ".pi/agent/models-store.json", ".pi/agent/settings.json", ".pi/settings.json"} {
			if err := copyIfPresent(filepath.Join(original, rel), filepath.Join(home, rel)); err != nil {
				return nil, err
			}
		}
	}
	env := map[string]string{}
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 && (baseEnvironment[parts[0]] || authEnvironment[parts[0]]) {
			env[parts[0]] = parts[1]
		}
	}
	if source := env["GOOGLE_APPLICATION_CREDENTIALS"]; source != "" {
		destination := filepath.Join(home, "credentials", "google-application-credentials.json")
		if err := copyIfPresent(source, destination); err != nil {
			return nil, err
		}
		env["GOOGLE_APPLICATION_CREDENTIALS"] = destination
	}
	env["HOME"] = home
	env["XDG_CONFIG_HOME"] = filepath.Join(home, ".config")
	env["XDG_CACHE_HOME"] = filepath.Join(home, ".cache")
	env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = "1"
	env["DISABLE_AUTOUPDATER"] = "1"
	env["PI_SKIP_VERSION_CHECK"] = "1"
	env["PI_TELEMETRY"] = "0"
	env["NO_COLOR"] = "1"
	return env, nil
}
func envList(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		out = append(out, k+"="+values[k])
	}
	return out
}

func makeSandbox(base, adapter string, executables, runtimeRoots []string) (sandboxContext, error) {
	if canonical, err := filepath.EvalSymlinks(filepath.Dir(base)); err == nil {
		base = filepath.Join(canonical, filepath.Base(base))
	}
	root := filepath.Join(base, "sandbox")
	home := filepath.Join(root, "home")
	temp := filepath.Join(root, "tmp")
	cwd := filepath.Join(root, "cwd")
	for _, p := range []string{temp, cwd} {
		if err := os.MkdirAll(p, 0o700); err != nil {
			return sandboxContext{}, err
		}
	}
	env, err := prepareEphemeralHome(home, adapter)
	if err != nil {
		return sandboxContext{}, err
	}
	env["TMPDIR"] = temp
	env["PATH"] = "/usr/bin:/bin"
	profile := filepath.Join(root, "profile.sb")
	if err := atomicWrite(profile, []byte(profileText(root, executables, runtimeRoots, adapter == "agy")), 0o600); err != nil {
		return sandboxContext{}, err
	}
	return sandboxContext{Root: root, Home: home, Temp: temp, CWD: cwd, Profile: profile, Environment: envList(env)}, nil
}

func currentUserCheck() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	if u.Uid == "0" {
		return fmt.Errorf("refusing to run experiments as root")
	}
	if required := os.Getenv("COUNCIL_EXPERIMENT_USER"); required != "" && required != u.Username {
		return fmt.Errorf("COUNCIL_EXPERIMENT_USER requires %q, current user is %q", required, u.Username)
	}
	return nil
}

func (h *Harness) sandboxProbe(verbose bool) error {
	if err := currentUserCheck(); err != nil {
		return err
	}
	base, err := os.MkdirTemp("", "council-seatbelt-probe-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(base)
	sentinel := filepath.Join(h.Root, "README.md")
	sandbox, err := makeSandbox(base, "mock", []string{"/bin/sh"}, nil)
	if err != nil {
		return err
	}
	script := `set -eu
printf allowed > "$TMPDIR/allowed"
if IFS= read -r line < "$DENIED_SENTINEL"; then exit 91; fi
exit 0`
	env := append(sandbox.Environment, "DENIED_SENTINEL="+sentinel)
	cmd := exec.Command("/usr/bin/sandbox-exec", "-f", sandbox.Profile, "--", "/bin/sh", "-c", script)
	cmd.Dir = sandbox.CWD
	cmd.Env = env
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Seatbelt probe failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	if data, err := os.ReadFile(filepath.Join(sandbox.Temp, "allowed")); err != nil || string(data) != "allowed" {
		return fmt.Errorf("Seatbelt did not permit isolated scratch write")
	}
	if verbose {
		fmt.Println("Seatbelt probe passed: scratch write allowed; repository read denied.")
	}
	return nil
}
func (h *Harness) SandboxProbe() error { return h.sandboxProbe(true) }
