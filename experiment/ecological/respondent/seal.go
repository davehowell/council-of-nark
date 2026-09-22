package respondent

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/davehowell/council-of-nark/experiment/ecological/snapshot"
)

type Seal struct {
	SchemaVersion int        `json:"schema_version"`
	CompletedAt   string     `json:"completed_at"`
	TreeSHA256    string     `json:"tree_sha256"`
	Files         []SealFile `json:"files"`
}

type SealFile struct {
	Path   string `json:"path"`
	Bytes  int64  `json:"bytes"`
	SHA256 string `json:"sha256"`
}

// RemoveEphemeralRunState deletes copied credentials and reproducible build
// caches before the immutable raw-artifact seal is created.
func RemoveEphemeralRunState(root string) error {
	for _, relative := range []string{"pi-sandbox/home", "pi-sandbox/tmp", "pi-sandbox/cwd"} {
		if err := os.RemoveAll(filepath.Join(root, filepath.FromSlash(relative))); err != nil {
			return err
		}
	}
	matches, err := filepath.Glob(filepath.Join(root, "focused-tests", "*", "scratch"))
	if err != nil {
		return err
	}
	for _, match := range matches {
		if err := os.RemoveAll(match); err != nil {
			return err
		}
	}
	return nil
}

func SealAttempt(root string) (Seal, error) {
	if _, err := os.Stat(filepath.Join(root, "seal.json")); err == nil || !os.IsNotExist(err) {
		return Seal{}, fmt.Errorf("attempt already has a seal")
	}
	files, err := sealedFiles(root)
	if err != nil {
		return Seal{}, err
	}
	seal := Seal{SchemaVersion: 1, CompletedAt: time.Now().UTC().Format(time.RFC3339Nano), Files: files}
	seal.TreeSHA256 = sealTreeDigest(files)
	if err := writeJSON(filepath.Join(root, "seal.json"), seal); err != nil {
		return Seal{}, err
	}
	for _, file := range files {
		if err := os.Chmod(filepath.Join(root, filepath.FromSlash(file.Path)), 0o400); err != nil {
			return Seal{}, err
		}
	}
	if err := os.Chmod(filepath.Join(root, "seal.json"), 0o400); err != nil {
		return Seal{}, err
	}
	return seal, nil
}

func VerifyAttemptSeal(root string) (Seal, error) {
	data, err := os.ReadFile(filepath.Join(root, "seal.json"))
	if err != nil {
		return Seal{}, err
	}
	var expected Seal
	if err := json.Unmarshal(data, &expected); err != nil {
		return Seal{}, err
	}
	if expected.SchemaVersion != 1 {
		return Seal{}, fmt.Errorf("unsupported seal schema_version %d", expected.SchemaVersion)
	}
	actual, err := sealedFiles(root)
	if err != nil {
		return Seal{}, err
	}
	if len(actual) != len(expected.Files) {
		return Seal{}, fmt.Errorf("sealed file count changed: got %d, expected %d", len(actual), len(expected.Files))
	}
	for index := range actual {
		if actual[index] != expected.Files[index] {
			return Seal{}, fmt.Errorf("sealed file changed: %s", actual[index].Path)
		}
	}
	if digest := sealTreeDigest(actual); digest != expected.TreeSHA256 {
		return Seal{}, fmt.Errorf("sealed tree digest changed")
	}
	return expected, nil
}

func sealedFiles(root string) ([]SealFile, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	files := []SealFile{}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root || entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("run artifact contains a symbolic link: %s", path)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("run artifact is not a regular file: %s", path)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == "seal.json" {
			return nil
		}
		lower := strings.ToLower(relative)
		if strings.HasSuffix(lower, "/auth.json") || strings.HasSuffix(lower, "/models-store.json") || strings.Contains(lower, "google-application-credentials") {
			return fmt.Errorf("credential material remained in run artifacts: %s", relative)
		}
		digest, err := snapshot.FileSHA256(path)
		if err != nil {
			return err
		}
		files = append(files, SealFile{Path: relative, Bytes: info.Size(), SHA256: digest})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func sealTreeDigest(files []SealFile) string {
	hash := sha256.New()
	for _, file := range files {
		fmt.Fprintf(hash, "%s\x00%d\x00%s\n", file.Path, file.Bytes, file.SHA256)
	}
	return hex.EncodeToString(hash.Sum(nil))
}
