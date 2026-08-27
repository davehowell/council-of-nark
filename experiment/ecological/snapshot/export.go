package snapshot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func (r *Runner) gitCache() string {
	return filepath.Join(r.Cache, "git", r.Config.TaskID+".git")
}

func (r *Runner) gitArgs(args ...string) []string {
	return append([]string{"--git-dir", r.gitCache()}, args...)
}

func (r *Runner) gitCapture(args ...string) (string, error) {
	return capture(r.Root, nil, "git", r.gitArgs(args...)...)
}

func (r *Runner) fetchAndVerify() error {
	cache := r.gitCache()
	if _, err := os.Stat(cache); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(cache), 0o700); err != nil {
			return err
		}
		result := runCommand(r.Root, nil, "git", "init", "--bare", cache)
		if err := writeCommandLog(filepath.Join(r.Controller, "logs", "git-init.json"), []string{"git", "init", "--bare", "<controller-cache>"}, result); err != nil {
			return err
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("git init failed: %s", strings.TrimSpace(result.Stderr))
		}
	} else if err != nil {
		return err
	}
	for index, commit := range []string{r.Config.ParentCommit, r.Config.EvidenceCommit} {
		args := r.gitArgs("fetch", "--no-tags", "--force", r.Config.RepositoryURL, commit)
		result := runCommand(r.Root, nil, "git", args...)
		logName := fmt.Sprintf("git-fetch-%d.json", index+1)
		if err := writeCommandLog(filepath.Join(r.Controller, "logs", logName), append([]string{"git"}, args...), result); err != nil {
			return err
		}
		if result.ExitCode != 0 {
			return fmt.Errorf("fetch exact commit %s failed: %s", commit, strings.TrimSpace(result.Stderr))
		}
	}
	for label, commit := range map[string]string{"parent": r.Config.ParentCommit, "evidence": r.Config.EvidenceCommit} {
		resolved, err := r.gitCapture("rev-parse", commit+"^{commit}")
		if err != nil {
			return err
		}
		if strings.TrimSpace(resolved) != commit {
			return fmt.Errorf("%s commit resolved to %s", label, strings.TrimSpace(resolved))
		}
	}
	parentsText, err := r.gitCapture("show", "-s", "--format=%P", r.Config.EvidenceCommit)
	if err != nil {
		return err
	}
	parents := strings.Fields(parentsText)
	if len(parents) == 0 || parents[0] != r.Config.ExpectedEvidenceParent {
		return fmt.Errorf("evidence first parent is %v, expected %s", parents, r.Config.ExpectedEvidenceParent)
	}
	fsck := runCommand(r.Root, nil, "git", r.gitArgs("fsck", "--connectivity-only", "--no-reflogs", r.Config.ParentCommit, r.Config.EvidenceCommit)...)
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", "git-fsck.json"), append([]string{"git"}, r.gitArgs("fsck", "--connectivity-only", "--no-reflogs", r.Config.ParentCommit, r.Config.EvidenceCommit)...), fsck); err != nil {
		return err
	}
	if fsck.ExitCode != 0 {
		return fmt.Errorf("git fsck failed: %s", strings.TrimSpace(fsck.Stderr))
	}
	metadata := map[string]any{}
	for label, commit := range map[string]string{"parent": r.Config.ParentCommit, "evidence": r.Config.EvidenceCommit} {
		value, err := r.gitCapture("show", "-s", "--format=%H%x00%T%x00%P%x00%aI%x00%cI", commit)
		if err != nil {
			return err
		}
		fields := strings.Split(strings.TrimSpace(value), "\x00")
		if len(fields) != 5 {
			return fmt.Errorf("unexpected commit metadata for %s", label)
		}
		metadata[label] = map[string]any{
			"commit": fields[0], "tree": fields[1], "parents": strings.Fields(fields[2]),
			"author_date": fields[3], "committer_date": fields[4],
		}
	}
	metadata["repository_url"] = r.Config.RepositoryURL
	r.provenance.Upstream = metadata

	for _, expected := range r.Config.LicenseFiles {
		data, err := r.gitBlob(r.Config.ParentCommit, expected.Path)
		if err != nil {
			return fmt.Errorf("read license %s: %w", expected.Path, err)
		}
		digest := shaBytes(data)
		if digest != expected.SHA256 {
			return fmt.Errorf("license digest changed for %s: got %s", expected.Path, digest)
		}
		if !bytes.Contains(data, []byte(expected.Contains)) {
			return fmt.Errorf("license marker missing from %s", expected.Path)
		}
		r.provenance.Licenses = append(r.provenance.Licenses, map[string]any{
			"path": expected.Path, "sha256": digest, "spdx": r.Config.LicenseSPDX,
		})
	}
	return nil
}

func (r *Runner) gitBlob(commit, path string) ([]byte, error) {
	command := exec.Command("git", r.gitArgs("show", commit+":"+path)...)
	command.Dir = r.Root
	var stderr bytes.Buffer
	command.Stderr = &stderr
	out, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("git show %s:%s: %w: %s", commit, path, err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func (r *Runner) writeArchive(commit, path, logName string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	command := exec.Command("git", r.gitArgs("archive", "--format=tar", commit)...)
	command.Dir = r.Root
	var stderr bytes.Buffer
	command.Stdout = file
	command.Stderr = &stderr
	runErr := command.Run()
	closeErr := file.Close()
	result := commandResult{ExitCode: 0, Stderr: stderr.String()}
	if runErr != nil {
		result.ExitCode = -1
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		}
	}
	if err := writeCommandLog(filepath.Join(r.Controller, "logs", logName), append([]string{"git"}, r.gitArgs("archive", "--format=tar", commit)...), result); err != nil {
		return err
	}
	if runErr != nil {
		return fmt.Errorf("git archive %s: %w: %s", commit, runErr, strings.TrimSpace(stderr.String()))
	}
	return closeErr
}

func (r *Runner) archiveParent() error {
	archive := filepath.Join(r.Controller, "archives", "parent.tar")
	if err := r.writeArchive(r.Config.ParentCommit, archive, "git-archive-parent.json"); err != nil {
		return err
	}
	digest, err := shaFile(archive)
	if err != nil {
		return err
	}
	stat, err := os.Stat(archive)
	if err != nil {
		return err
	}
	r.provenance.Archive = map[string]any{"format": "git archive --format=tar", "sha256": digest, "bytes": stat.Size()}
	return extractTar(archive, r.Source)
}

func (r *Runner) sanitizeSource() error {
	for _, removal := range r.Config.Sanitization.Remove {
		path := filepath.Join(r.Source, filepath.FromSlash(removal.Path))
		if _, err := os.Lstat(path); err != nil {
			if os.IsNotExist(err) && !removal.Required {
				continue
			}
			if os.IsNotExist(err) {
				return fmt.Errorf("required sanitization path is absent: %s", removal.Path)
			}
			return err
		}
		manifest, err := manifestOne(r.Source, removal.Path)
		if err != nil {
			return err
		}
		r.provenance.Sanitization = append(r.provenance.Sanitization, map[string]any{
			"path": removal.Path, "reason": removal.Reason, "required": removal.Required,
			"removed_tree_sha256": manifest.TreeSHA256, "entries": manifest.FileCount, "bytes": manifest.ByteCount,
		})
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	if err := scanForbiddenPaths(r.Source); err != nil {
		return err
	}
	if err := scanForbiddenContent(r.Source, r.Config.Sanitization.ForbiddenContent); err != nil {
		return err
	}
	if _, err := os.Lstat(filepath.Join(r.Source, ".git")); !os.IsNotExist(err) {
		return fmt.Errorf("sanitized source contains .git")
	}
	for _, expected := range r.Config.LicenseFiles {
		path := filepath.Join(r.Source, filepath.FromSlash(expected.Path))
		digest, err := shaFile(path)
		if err != nil {
			return err
		}
		if digest != expected.SHA256 {
			return fmt.Errorf("retained license %s changed during sanitization", expected.Path)
		}
	}
	if err := normalizeSource(r.Source, false); err != nil {
		return err
	}
	manifest, err := treeManifest(r.Source)
	if err != nil {
		return err
	}
	if err := writeJSON(filepath.Join(r.Controller, "source-manifest-preliminary.json"), manifest, 0o600); err != nil {
		return err
	}
	r.provenance.Source = map[string]any{
		"history": "none; .git is absent", "mtime": neutralTime.Format("2006-01-02T15:04:05Z"),
		"preliminary_tree_sha256": manifest.TreeSHA256,
	}
	return nil
}

func manifestOne(root, relative string) (TreeManifest, error) {
	path := filepath.Join(root, filepath.FromSlash(relative))
	info, err := os.Lstat(path)
	if err != nil {
		return TreeManifest{}, err
	}
	if info.IsDir() {
		return treeManifest(path)
	}
	digest := ""
	target := ""
	typeName := "file"
	if info.Mode().IsRegular() {
		digest, err = shaFile(path)
	} else if info.Mode()&os.ModeSymlink != 0 {
		typeName = "symlink"
		target, err = os.Readlink(path)
	} else {
		return TreeManifest{}, fmt.Errorf("unsupported removal entry %s", relative)
	}
	if err != nil {
		return TreeManifest{}, err
	}
	record := FileRecord{Path: relative, Type: typeName, Mode: uint32(info.Mode().Perm()), Size: info.Size(), SHA256: digest, Target: target}
	canonical, _ := jsonMarshal([]FileRecord{record})
	return TreeManifest{SchemaVersion: 1, TreeSHA256: shaBytes(canonical), FileCount: 1, ByteCount: info.Size(), Files: []FileRecord{record}}, nil
}

func jsonMarshal(value any) ([]byte, error) {
	// Kept as a small seam so manifestOne and treeManifest use the same JSON encoding.
	return json.Marshal(value)
}

var forbiddenSegments = map[string]bool{
	".git": true, ".github": true, ".cursor": true, ".claude": true, ".agents": true,
	".mcp": true, ".continue": true, ".gemini": true, ".kiro": true, ".vscode": true,
	"patches": true,
}
var forbiddenBaseNames = map[string]bool{
	"agents.md": true, "agent.md": true, "claude.md": true, "gemini.md": true,
	"copilot.md": true, "copilot-instructions.md": true, ".cursorrules": true,
	".windsurfrules": true, ".mcp.json": true, "mcp.json": true, "opencode.json": true,
}

func scanForbiddenPaths(root string) error {
	var found []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		parts := strings.Split(rel, "/")
		for _, part := range parts {
			if forbiddenSegments[strings.ToLower(part)] {
				found = append(found, rel)
				if entry.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}
		base := strings.ToLower(parts[len(parts)-1])
		if forbiddenBaseNames[base] || strings.HasSuffix(base, ".patch") || strings.HasSuffix(base, ".diff") ||
			strings.HasPrefix(base, "changelog") || strings.HasPrefix(base, "changes.") || strings.HasPrefix(base, "history.") || strings.HasPrefix(base, "news.") {
			found = append(found, rel)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(found) > 0 {
		sort.Strings(found)
		return fmt.Errorf("unconfigured forbidden paths remain after sanitization: %s", strings.Join(found, ", "))
	}
	return nil
}

func scanForbiddenContent(root string, needles []string) error {
	var hits []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, needle := range needles {
			if bytes.Contains(data, []byte(needle)) {
				rel, _ := filepath.Rel(root, path)
				hits = append(hits, filepath.ToSlash(rel)+" (configured evidence marker)")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if len(hits) > 0 {
		sort.Strings(hits)
		return fmt.Errorf("forbidden evidence content remains: %s", strings.Join(hits, ", "))
	}
	return nil
}

func copyBlobTo(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	_, writeErr := io.Copy(file, bytes.NewReader(data))
	closeErr := file.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
