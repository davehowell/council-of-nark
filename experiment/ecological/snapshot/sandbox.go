package snapshot

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func seatbeltQuote(value string) string { return strconv.Quote(value) }

func uniquePaths(values []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func ecologicalProfile(readRoots, executableFiles, executableRoots []string, scratch string) string {
	readRoots = append(readRoots,
		"/System", "/usr/lib", "/usr/share", "/Library/Apple", "/private/etc",
		"/private/var/db/timezone", "/dev", scratch,
	)
	var builder strings.Builder
	builder.WriteString("(version 1)\n(deny default)\n")
	builder.WriteString("(allow file-read-metadata)\n")
	builder.WriteString("(allow sysctl-read)\n(allow mach-lookup)\n(allow system-socket)\n(allow process-fork)\n")
	builder.WriteString("(allow file-read*\n  (literal \"/\")\n")
	for _, path := range uniquePaths(readRoots) {
		builder.WriteString("  (subpath " + seatbeltQuote(path) + ")\n")
	}
	for _, path := range uniquePaths(executableFiles) {
		builder.WriteString("  (literal " + seatbeltQuote(path) + ")\n")
	}
	builder.WriteString(")\n")
	builder.WriteString("(allow file-write* (subpath " + seatbeltQuote(scratch) + ") (literal \"/dev/null\"))\n")
	builder.WriteString("(allow process-exec\n")
	for _, path := range uniquePaths(executableFiles) {
		builder.WriteString("  (literal " + seatbeltQuote(path) + ")\n")
	}
	for _, path := range uniquePaths(executableRoots) {
		builder.WriteString("  (subpath " + seatbeltQuote(path) + ")\n")
	}
	builder.WriteString(")\n")
	return builder.String()
}

func denyNetworkProbe(profile, cwd string, env []string) (commandResult, error) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return commandResult{}, err
	}
	defer listener.Close()
	if tcp, ok := listener.(*net.TCPListener); ok {
		_ = tcp.SetDeadline(time.Now().Add(3 * time.Second))
	}
	address := listener.Addr().(*net.TCPAddr)
	args := []string{"-f", profile, "--", "/usr/bin/nc", "-z", "-w", "1", "127.0.0.1", strconv.Itoa(address.Port)}
	result := runCommand(cwd, env, "/usr/bin/sandbox-exec", args...)
	accepted := make(chan bool, 1)
	go func() {
		connection, acceptErr := listener.Accept()
		if acceptErr == nil {
			connection.Close()
			accepted <- true
			return
		}
		accepted <- false
	}()
	connected := <-accepted
	if result.ExitCode == 0 || connected {
		return result, fmt.Errorf("Seatbelt network-denial probe reached a controller listener")
	}
	return result, nil
}

func (r *Runner) runOfflineRegression() error {
	var closure struct {
		ToolchainPath string         `json:"toolchain_path"`
		ToolchainTree string         `json:"toolchain_tree_sha256"`
		ModulePath    string         `json:"module_path"`
		ModuleTree    string         `json:"module_tree_sha256"`
		Modules       []ModuleRecord `json:"modules"`
	}
	if err := readJSON(filepath.Join(r.Controller, "closure.json"), &closure); err != nil {
		return err
	}
	toolchain := filepath.Join(r.Ecological, filepath.FromSlash(closure.ToolchainPath))
	modules := filepath.Join(r.Ecological, filepath.FromSlash(closure.ModulePath))
	paths := closurePaths{
		Toolchain: toolchain, ToolchainDigest: closure.ToolchainTree,
		Modules: modules, ModulesDigest: closure.ModuleTree,
		GoBinary: filepath.Join(toolchain, "root", "bin", "go"), ModuleRecords: closure.Modules,
	}
	if err := r.verifyClosure(paths); err != nil {
		return err
	}

	profileRoot := filepath.Join(r.Controller, "offline-sandbox")
	scratch := filepath.Join(profileRoot, "scratch")
	if err := os.MkdirAll(scratch, 0o700); err != nil {
		return err
	}
	readRoots := []string{
		filepath.Join(r.Controller, "validation", "parent"), filepath.Join(r.Controller, "validation", "evidence"),
		filepath.Join(toolchain, "root"), filepath.Join(modules, "gomod"),
	}
	executableFiles := []string{paths.GoBinary, "/usr/bin/nc"}
	executableRoots := []string{filepath.Join(toolchain, "root", "pkg", "tool"), scratch}
	profileText := ecologicalProfile(readRoots, executableFiles, executableRoots, scratch)
	profilePath := filepath.Join(profileRoot, "profile.sb")
	if err := atomicWrite(profilePath, []byte(profileText), 0o600); err != nil {
		return err
	}
	env := goEnvironment(paths, filepath.Join(modules, "gomod"), filepath.Join(scratch, "gocache"), false)
	networkResult, err := denyNetworkProbe(profilePath, scratch, env)
	if logErr := writeCommandLog(filepath.Join(r.Controller, "logs", "offline-network-probe.json"), []string{"sandbox-exec", "<offline-profile>", "nc", "<controller-listener>"}, networkResult); logErr != nil {
		return logErr
	}
	if err != nil {
		return err
	}
	executableResult := runSandboxed(profilePath, scratch, env, "/bin/sh", "-c", "true")
	if logErr := writeCommandLog(filepath.Join(r.Controller, "logs", "offline-executable-probe.json"), []string{"sandbox-exec", "<offline-profile>", "sh", "<probe>"}, executableResult); logErr != nil {
		return logErr
	}
	if executableResult.ExitCode == 0 {
		return fmt.Errorf("offline Seatbelt profile executed an unlisted shell")
	}

	evidenceRoot := filepath.Join(r.Controller, "validation", "evidence")
	listArgs := []string{"list", "-deps", "-test", "-json", r.Config.FocusedTest.Package}
	listResult := runSandboxed(profilePath, evidenceRoot, env, paths.GoBinary, listArgs...)
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "offline-go-list.json"), append([]string{"<frozen-go>"}, listArgs...), listResult); err != nil {
		return err
	}
	if listResult.ExitCode != 0 {
		return fmt.Errorf("offline go list failed: %s", strings.TrimSpace(listResult.Stderr))
	}
	offlineModules, err := parseModules([]byte(listResult.Stdout), filepath.Join(modules, "gomod"))
	if err != nil {
		return err
	}
	if !equalModuleIdentity(paths.ModuleRecords, offlineModules) {
		return fmt.Errorf("offline module identities differ from the prefetched closure")
	}
	if err := atomicWrite(filepath.Join(r.Controller, "dependency-packages-offline.jsonstream"), []byte(listResult.Stdout), 0o600); err != nil {
		return err
	}

	verifyResult := runSandboxed(profilePath, evidenceRoot, env, paths.GoBinary, "mod", "verify")
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "offline-go-mod-verify.json"), []string{"<frozen-go>", "mod", "verify"}, verifyResult); err != nil {
		return err
	}
	if verifyResult.ExitCode != 0 {
		return fmt.Errorf("offline go mod verify failed: %s", strings.TrimSpace(verifyResult.Stderr))
	}

	testArgs := []string{"test", "-count=1", r.Config.FocusedTest.Package, "-run", r.Config.FocusedTest.Run}
	parentRoot := filepath.Join(r.Controller, "validation", "parent")
	parentResult := runSandboxed(profilePath, parentRoot, env, paths.GoBinary, testArgs...)
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "parent-regression.json"), append([]string{"<frozen-go>"}, testArgs...), parentResult); err != nil {
		return err
	}
	parentOutput := parentResult.Stdout + "\n" + parentResult.Stderr
	if parentResult.ExitCode == 0 || !strings.Contains(parentOutput, r.Config.FocusedTest.ParentFailureContains) {
		return fmt.Errorf("parent did not produce the expected discriminating failure (exit %d)", parentResult.ExitCode)
	}
	evidenceResult := runSandboxed(profilePath, evidenceRoot, env, paths.GoBinary, testArgs...)
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "evidence-regression.json"), append([]string{"<frozen-go>"}, testArgs...), evidenceResult); err != nil {
		return err
	}
	if evidenceResult.ExitCode != 0 {
		return fmt.Errorf("evidence regression test failed: %s", strings.TrimSpace(evidenceResult.Stdout+"\n"+evidenceResult.Stderr))
	}
	r.provenance.Validation["parent"] = map[string]any{"exit_code": parentResult.ExitCode, "expected_failure_observed": true}
	r.provenance.Validation["evidence"] = map[string]any{"exit_code": evidenceResult.ExitCode, "passed": true}
	r.provenance.Validation["command"] = append([]string{"go"}, testArgs...)
	r.provenance.Validation["closure_verified_offline"] = true
	r.provenance.Validation["network_probe_denied"] = true
	r.provenance.Validation["unlisted_executable_denied"] = true
	r.provenance.Validation["seatbelt_profile_sha256"] = shaBytes([]byte(profileText))
	return nil
}

func runSandboxed(profile, cwd string, env []string, executable string, args ...string) commandResult {
	sandboxArgs := []string{"-f", profile, "--", executable}
	sandboxArgs = append(sandboxArgs, args...)
	return runCommand(cwd, env, "/usr/bin/sandbox-exec", sandboxArgs...)
}

func (r *Runner) verifyClosure(paths closurePaths) error {
	toolchain, err := treeManifestExcludingMetadata(paths.Toolchain)
	if err != nil {
		return err
	}
	if toolchain.TreeSHA256 != paths.ToolchainDigest {
		return fmt.Errorf("frozen toolchain digest verification failed")
	}
	modules, err := treeManifestExcludingMetadata(paths.Modules)
	if err != nil {
		return err
	}
	if modules.TreeSHA256 != paths.ModulesDigest {
		return fmt.Errorf("frozen module digest verification failed")
	}
	return nil
}

func equalModuleIdentity(left, right []ModuleRecord) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index].Path != right[index].Path || left[index].Version != right[index].Version || left[index].Sum != right[index].Sum {
			return false
		}
	}
	return true
}

func (r *Runner) runSourceIsolationProbe() error {
	preliminary, err := treeManifest(r.Source)
	if err != nil {
		return err
	}
	probeRoot := filepath.Join(r.Controller, "source-probe")
	scratch := filepath.Join(probeRoot, "scratch")
	if err := os.MkdirAll(scratch, 0o700); err != nil {
		return err
	}
	sibling := filepath.Join(r.Controller, "sibling-task", "secret.txt")
	if err := atomicWrite(sibling, []byte("sibling task sentinel\n"), 0o600); err != nil {
		return err
	}
	evidence := filepath.Join(r.Controller, "validation", "evidence", filepath.FromSlash(r.Config.FocusedTest.EvidenceTestFiles[0]))
	metadata := filepath.Join(r.Controller, "config.json")
	council := filepath.Join(r.Root, "experiment", "ecological", "candidates.json")
	for _, path := range []string{evidence, metadata, council, sibling} {
		data, err := os.ReadFile(path)
		if err != nil || len(data) == 0 {
			return fmt.Errorf("isolation sentinel is unavailable: %s: %v", path, err)
		}
	}
	writeTarget := filepath.Join(r.Source, ".gitignore")
	beforeWrite, err := shaFile(writeTarget)
	if err != nil {
		return err
	}
	profileText := ecologicalProfile([]string{r.Source}, []string{"/bin/sh", "/bin/bash", "/usr/bin/nc"}, nil, scratch)
	profilePath := filepath.Join(probeRoot, "profile.sb")
	if err := atomicWrite(profilePath, []byte(profileText), 0o600); err != nil {
		return err
	}
	script := `set -eu
IFS= read -r first < "$SOURCE/LICENSE.md"
test -n "$first"
for denied in "$COUNCIL" "$METADATA" "$SIBLING" "$EVIDENCE" "$SOURCE/.git/HEAD"; do
  if IFS= read -r leaked < "$denied"; then exit 91; fi
done
if printf tamper >> "$WRITE_TARGET"; then exit 92; fi
exit 0`
	env := sortedEnvironment(map[string]string{
		"HOME": filepath.Join(scratch, "home"), "TMPDIR": filepath.Join(scratch, "tmp"),
		"PATH": "/usr/bin:/bin", "SOURCE": r.Source, "COUNCIL": council, "METADATA": metadata,
		"SIBLING": sibling, "EVIDENCE": evidence, "WRITE_TARGET": writeTarget, "LANG": "en_US.UTF-8",
	})
	_ = os.MkdirAll(filepath.Join(scratch, "home"), 0o700)
	_ = os.MkdirAll(filepath.Join(scratch, "tmp"), 0o700)
	result := runSandboxed(profilePath, scratch, env, "/bin/sh", "-c", script)
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "source-filesystem-probe.json"), []string{"sandbox-exec", "<source-profile>", "sh", "<probe>"}, result); err != nil {
		return err
	}
	if result.ExitCode != 0 {
		return fmt.Errorf("source filesystem isolation probe failed: %s", strings.TrimSpace(result.Stderr))
	}
	afterWrite, err := shaFile(writeTarget)
	if err != nil {
		return err
	}
	if beforeWrite != afterWrite {
		return fmt.Errorf("source write target changed during isolation probe")
	}
	networkResult, err := denyNetworkProbe(profilePath, scratch, env)
	if logErr := writeCommandLog(filepath.Join(r.Controller, "logs", "source-network-probe.json"), []string{"sandbox-exec", "<source-profile>", "nc", "<controller-listener>"}, networkResult); logErr != nil {
		return logErr
	}
	if err != nil {
		return err
	}
	after, err := treeManifest(r.Source)
	if err != nil {
		return err
	}
	if after.TreeSHA256 != preliminary.TreeSHA256 {
		return fmt.Errorf("source tree changed during isolation probes")
	}
	r.provenance.Isolation = map[string]any{
		"source_read": true, "source_write_denied": true, "council_read_denied": true,
		"controller_metadata_read_denied": true, "sibling_task_read_denied": true,
		"evidence_read_denied": true, "git_history_absent": true, "network_denied": true,
		"seatbelt_profile_sha256": shaBytes([]byte(profileText)),
	}
	return nil
}

func (r *Runner) finalize() error {
	if err := normalizeSource(r.Source, true); err != nil {
		return err
	}
	manifest, err := treeManifest(r.Source)
	if err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(r.Controller, "source-manifest.json"), manifest, 0o600); err != nil {
		return err
	}
	r.provenance.Source["tree_sha256"] = manifest.TreeSHA256
	r.provenance.Source["entries"] = manifest.FileCount
	r.provenance.Source["bytes"] = manifest.ByteCount
	r.provenance.Source["read_only"] = true
	r.provenance.Source["provenance_inside_source"] = false
	if err := writeJSON(filepath.Join(r.Controller, "provenance.json"), r.provenance, 0o600); err != nil {
		return err
	}
	provenanceDigest, err := shaFile(filepath.Join(r.Controller, "provenance.json"))
	if err != nil {
		return err
	}
	return writeJSON(filepath.Join(r.Controller, "seal.json"), map[string]any{
		"schema_version": 1, "provenance_sha256": provenanceDigest,
		"source_manifest_sha256": mustSHA(filepath.Join(r.Controller, "source-manifest.json")),
		"source_tree_sha256":     manifest.TreeSHA256,
	}, 0o600)
}

func mustSHA(path string) string {
	digest, _ := shaFile(path)
	return digest
}

func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
