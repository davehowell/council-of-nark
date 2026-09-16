package snapshot

import (
	"archive/tar"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func repositoryRoot(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(out))
}

func TestGortexConfigValidates(t *testing.T) {
	path := filepath.Join(repositoryRoot(t), "experiment/ecological/config/eco-gortex-unicode-tokenizer.json")
	config, _, err := loadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if config.TaskID != "eco-gortex-unicode-tokenizer" || config.FocusedTest.Package != "./internal/search/rerank" {
		t.Fatalf("unexpected config: %#v", config)
	}
}

func TestExtractTarRejectsTraversal(t *testing.T) {
	archive := filepath.Join(t.TempDir(), "bad.tar")
	file, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	writer := tar.NewWriter(file)
	if err := writer.WriteHeader(&tar.Header{Name: "../escape", Mode: 0o644, Size: 1, Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write([]byte("x")); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := extractTar(archive, filepath.Join(t.TempDir(), "out")); err == nil || !strings.Contains(err.Error(), "unsafe archive path") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestForbiddenPathScanFailsClosed(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("instructions"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := scanForbiddenPaths(root); err == nil || !strings.Contains(err.Error(), "CLAUDE.md") {
		t.Fatalf("unexpected scan result: %v", err)
	}
}

func TestTreeManifestDetectsContentChange(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "source.go")
	if err := os.WriteFile(path, []byte("package sample\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	first, err := treeManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := treeManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if first.TreeSHA256 != second.TreeSHA256 {
		t.Fatal("unchanged tree digest is not deterministic")
	}
	if err := os.WriteFile(path, []byte("package changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := treeManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if changed.TreeSHA256 == first.TreeSHA256 {
		t.Fatal("content change did not alter tree digest")
	}
}

func TestClosureManifestExcludesOnlyMetadata(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "root", "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "root", "bin", "go"), []byte("tool"), 0o555); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "notices"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "notices", "LICENSE"), []byte("license"), 0o444); err != nil {
		t.Fatal(err)
	}
	before, err := treeManifest(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte("metadata"), 0o444); err != nil {
		t.Fatal(err)
	}
	after, err := treeManifestExcludingMetadata(root)
	if err != nil {
		t.Fatal(err)
	}
	if before.TreeSHA256 != after.TreeSHA256 {
		t.Fatalf("closure digest changed: %s != %s", before.TreeSHA256, after.TreeSHA256)
	}
}

func TestRuntimeModuleBoundaryIsStable(t *testing.T) {
	root := filepath.Join(t.TempDir(), "runtime")
	if err := ensureRuntimeModuleBoundary(root, "council.local/test-runtime"); err != nil {
		t.Fatal(err)
	}
	if err := ensureRuntimeModuleBoundary(root, "council.local/test-runtime"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "module council.local/test-runtime\n\ngo 1.22\n" {
		t.Fatalf("unexpected boundary: %q", data)
	}
}

func TestEcologicalProfileHasNoNetworkAllowance(t *testing.T) {
	profile := ecologicalProfile([]string{"/private/tmp/source"}, []string{"/bin/sh"}, nil, "/private/tmp/scratch")
	if strings.Contains(profile, "network-outbound") || strings.Contains(profile, "network-inbound") || strings.Contains(profile, "(allow process*)") || !strings.Contains(profile, "(deny default)") {
		t.Fatalf("unsafe profile:\n%s", profile)
	}
	if !bytes.Contains([]byte(profile), []byte(`(subpath "/private/tmp/source")`)) {
		t.Fatal("source read root missing")
	}
}
