package snapshot

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileRecord struct {
	Path   string `json:"path"`
	Type   string `json:"type"`
	Mode   uint32 `json:"mode"`
	Size   int64  `json:"size,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
	Target string `json:"target,omitempty"`
}

type TreeManifest struct {
	SchemaVersion int          `json:"schema_version"`
	TreeSHA256    string       `json:"tree_sha256"`
	FileCount     int          `json:"file_count"`
	ByteCount     int64        `json:"byte_count"`
	Files         []FileRecord `json:"files"`
}

func shaBytes(value []byte) string {
	digest := sha256.Sum256(value)
	return hex.EncodeToString(digest[:])
}

// FileSHA256 returns the lowercase SHA-256 digest of one regular file.
func FileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	digest := sha256.New()
	if _, err := io.Copy(digest, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func shaFile(path string) (string, error) { return FileSHA256(path) }

// BuildTreeManifest hashes a source tree using the snapshot manifest format.
func BuildTreeManifest(root string) (TreeManifest, error) { return treeManifest(root) }

func treeManifest(root string) (TreeManifest, error) {
	var records []FileRecord
	var bytes int64
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
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		record := FileRecord{Path: rel, Mode: uint32(info.Mode().Perm())}
		switch {
		case info.Mode().IsRegular():
			record.Type = "file"
			record.Size = info.Size()
			record.SHA256, err = FileSHA256(path)
			bytes += info.Size()
		case info.IsDir():
			record.Type = "directory"
		case info.Mode()&os.ModeSymlink != 0:
			record.Type = "symlink"
			record.Target, err = os.Readlink(path)
		default:
			return fmt.Errorf("unsupported filesystem entry %s (%s)", rel, info.Mode())
		}
		if err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return TreeManifest{}, err
	}
	sort.Slice(records, func(i, j int) bool { return records[i].Path < records[j].Path })
	canonical, err := json.Marshal(records)
	if err != nil {
		return TreeManifest{}, err
	}
	return TreeManifest{SchemaVersion: 1, TreeSHA256: shaBytes(canonical), FileCount: len(records), ByteCount: bytes, Files: records}, nil
}

func extractTar(tarPath, destination string) error {
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	if err := os.Mkdir(destination, 0o755); err != nil {
		return err
	}
	file, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := tar.NewReader(file)
	var total int64
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Typeflag == tar.TypeXGlobalHeader || header.Typeflag == tar.TypeXHeader {
			continue
		}
		name := filepath.ToSlash(header.Name)
		if name == "" || strings.HasPrefix(name, "/") || strings.HasPrefix(name, "../") || filepath.ToSlash(filepath.Clean(name)) != strings.TrimSuffix(name, "/") {
			return fmt.Errorf("unsafe archive path %q", header.Name)
		}
		name = strings.TrimSuffix(name, "/")
		if name == "" {
			continue
		}
		target := filepath.Join(destination, filepath.FromSlash(name))
		rel, err := filepath.Rel(destination, target)
		if err != nil || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
			return fmt.Errorf("archive path escapes destination: %q", header.Name)
		}
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if header.Size < 0 || header.Size > 1<<30 {
				return fmt.Errorf("archive file %s has unreasonable size %d", name, header.Size)
			}
			total += header.Size
			if total > 2<<30 {
				return fmt.Errorf("archive exceeds 2 GiB extraction limit")
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			mode := os.FileMode(0o644)
			if header.FileInfo().Mode()&0o111 != 0 {
				mode = 0o755
			}
			out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
			if err != nil {
				return err
			}
			_, copyErr := io.CopyN(out, reader, header.Size)
			closeErr := out.Close()
			if copyErr != nil {
				return copyErr
			}
			if closeErr != nil {
				return closeErr
			}
		case tar.TypeSymlink:
			if filepath.IsAbs(header.Linkname) {
				return fmt.Errorf("absolute symlink in archive: %s", name)
			}
			resolved := filepath.Clean(filepath.Join(filepath.Dir(target), header.Linkname))
			rel, err := filepath.Rel(destination, resolved)
			if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				return fmt.Errorf("symlink escapes archive root: %s -> %s", name, header.Linkname)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported archive entry %s (type %d)", name, header.Typeflag)
		}
	}
	return nil
}

var neutralTime = time.Unix(946684800, 0).UTC()

func normalizeSource(root string, readOnly bool) error {
	var directories []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		if info.IsDir() {
			directories = append(directories, path)
			return nil
		}
		mode := os.FileMode(0o644)
		if info.Mode()&0o111 != 0 {
			mode = 0o755
		}
		if readOnly {
			mode &^= 0o222
		}
		if err := os.Chmod(path, mode); err != nil {
			return err
		}
		return os.Chtimes(path, neutralTime, neutralTime)
	})
	if err != nil {
		return err
	}
	for i := len(directories) - 1; i >= 0; i-- {
		mode := os.FileMode(0o755)
		if readOnly {
			mode = 0o555
		}
		if err := os.Chmod(directories[i], mode); err != nil {
			return err
		}
		if err := os.Chtimes(directories[i], neutralTime, neutralTime); err != nil {
			return err
		}
	}
	return nil
}

func makeTreeWritable(root string) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return os.Chmod(path, 0o755)
		}
		return nil
	})
}

func copyTree(source, destination string, readOnly bool) error {
	if err := os.Mkdir(destination, 0o755); err != nil {
		return err
	}
	return filepath.WalkDir(source, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == source {
			return nil
		}
		rel, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destination, rel)
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		switch {
		case info.IsDir():
			return os.Mkdir(target, 0o755)
		case info.Mode().IsRegular():
			mode := info.Mode().Perm()
			if readOnly {
				mode &^= 0o222
			}
			in, err := os.Open(path)
			if err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
			if err != nil {
				in.Close()
				return err
			}
			_, copyErr := io.Copy(out, in)
			inErr := in.Close()
			outErr := out.Close()
			if copyErr != nil {
				return copyErr
			}
			if inErr != nil {
				return inErr
			}
			return outErr
		case info.Mode()&os.ModeSymlink != 0:
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, target)
		default:
			return fmt.Errorf("unsupported entry while copying %s", path)
		}
	})
}
