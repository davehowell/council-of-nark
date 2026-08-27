package mediator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type listArguments struct {
	Path  string `json:"path"`
	Depth int    `json:"depth"`
}

type readArguments struct {
	Path      string `json:"path"`
	StartLine int    `json:"start_line"`
	LineCount int    `json:"line_count"`
}

type searchArguments struct {
	Pattern string `json:"pattern"`
	Path    string `json:"path"`
	Include string `json:"include"`
}

type testArguments struct {
	Target string `json:"target"`
}

func (s *Session) list(args listArguments) (string, map[string]any, error) {
	if args.Depth == 0 {
		args.Depth = 1
	}
	if args.Depth < 1 || args.Depth > s.Policy.MaxListDepth {
		return "", nil, fmt.Errorf("depth must be between 1 and %d", s.Policy.MaxListDepth)
	}
	root, relative, err := s.resolve(args.Path, true)
	if err != nil {
		return "", nil, err
	}
	type row struct {
		path string
		line string
	}
	rows := []row{}
	truncated := false
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		relToRoot, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		depth := strings.Count(filepath.ToSlash(relToRoot), "/") + 1
		if entry.IsDir() && depth > args.Depth {
			return filepath.SkipDir
		}
		if depth > args.Depth {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("source contains a symbolic link: %s", filepath.ToSlash(relToRoot))
		}
		if len(rows) >= s.Policy.MaxListEntries {
			truncated = true
			return filepath.SkipAll
		}
		fromSource, err := filepath.Rel(s.Source, path)
		if err != nil {
			return err
		}
		fromSource = filepath.ToSlash(fromSource)
		if entry.IsDir() {
			rows = append(rows, row{fromSource, "d " + fromSource + "/"})
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		rows = append(rows, row{fromSource, fmt.Sprintf("f %s %d", fromSource, info.Size())})
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].path < rows[j].path })
	lines := make([]string, len(rows))
	for index, row := range rows {
		lines[index] = row.line
	}
	if truncated {
		lines = append(lines, fmt.Sprintf("[entry limit reached: %d]", s.Policy.MaxListEntries))
	}
	return strings.Join(lines, "\n"), map[string]any{
		"path": relative, "depth": args.Depth, "entries": len(rows), "entry_limit_reached": truncated,
	}, nil
}

func (s *Session) read(args readArguments) (string, map[string]any, error) {
	if args.Path == "" {
		return "", nil, fmt.Errorf("path is required")
	}
	if args.StartLine == 0 {
		args.StartLine = 1
	}
	if args.LineCount == 0 {
		args.LineCount = s.Policy.MaxReadLines
	}
	if args.StartLine < 1 || args.LineCount < 1 || args.LineCount > s.Policy.MaxReadLines {
		return "", nil, fmt.Errorf("start_line must be positive and line_count must be between 1 and %d", s.Policy.MaxReadLines)
	}
	path, relative, err := s.resolve(args.Path, false)
	if err != nil {
		return "", nil, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", nil, err
	}
	if info.Size() > s.Policy.MaxTextFileBytes {
		return "", nil, fmt.Errorf("file exceeds text-file size limit")
	}
	lines := []string{}
	totalLines := 0
	text, err := scanTextFile(path, s.Policy.MaxTextFileBytes, func(line int, value string) bool {
		totalLines = line
		if line >= args.StartLine && line < args.StartLine+args.LineCount {
			lines = append(lines, fmt.Sprintf("%d:%s", line, value))
		}
		return line < args.StartLine+args.LineCount
	})
	if err != nil {
		return "", nil, err
	}
	if !text {
		return "", nil, fmt.Errorf("file is binary or is not readable as UTF-8 text")
	}
	if len(lines) == 0 {
		return "", nil, fmt.Errorf("start_line is beyond the end of the file")
	}
	return strings.Join(lines, "\n"), map[string]any{
		"path": relative, "start_line": args.StartLine, "lines_returned": len(lines), "last_line_scanned": totalLines,
	}, nil
}

func (s *Session) search(args searchArguments) (string, map[string]any, error) {
	if args.Pattern == "" {
		return "", nil, fmt.Errorf("pattern is required")
	}
	if len([]byte(args.Pattern)) > s.Policy.MaxSearchPattern {
		return "", nil, fmt.Errorf("pattern exceeds %d bytes", s.Policy.MaxSearchPattern)
	}
	pattern, err := regexp.Compile(args.Pattern)
	if err != nil {
		return "", nil, fmt.Errorf("invalid RE2 pattern: %w", err)
	}
	root, relative, err := s.resolve(args.Path, true)
	if err != nil {
		return "", nil, err
	}
	if args.Include != "" {
		if filepath.IsAbs(args.Include) || strings.Contains(args.Include, "..") || strings.ContainsAny(args.Include, `/\\`) {
			return "", nil, fmt.Errorf("include must be a basename glob such as *.go")
		}
		if _, err := filepath.Match(args.Include, "probe"); err != nil {
			return "", nil, fmt.Errorf("invalid include glob: %w", err)
		}
	}
	matches := []string{}
	filesScanned := 0
	filesSkipped := 0
	limitReached := false
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if limitReached {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("source contains a symbolic link")
		}
		if args.Include != "" {
			included, err := filepath.Match(args.Include, entry.Name())
			if err != nil || !included {
				return err
			}
		}
		text, err := scanTextFile(path, s.Policy.MaxTextFileBytes, func(line int, value string) bool {
			if !pattern.MatchString(value) {
				return true
			}
			rel, _ := filepath.Rel(s.Source, path)
			matches = append(matches, fmt.Sprintf("%s:%d:%s", filepath.ToSlash(rel), line, linePreview(value, 500)))
			if len(matches) >= s.Policy.MaxSearchMatches {
				limitReached = true
				return false
			}
			return true
		})
		if err != nil {
			return err
		}
		if text {
			filesScanned++
		} else {
			filesSkipped++
		}
		if limitReached {
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	if len(matches) == 0 {
		return "No matches found", map[string]any{
			"path": relative, "pattern": args.Pattern, "include": args.Include,
			"matches": 0, "files_scanned": filesScanned, "files_skipped": filesSkipped, "match_limit_reached": false,
		}, nil
	}
	if limitReached {
		matches = append(matches, fmt.Sprintf("[match limit reached: %d]", s.Policy.MaxSearchMatches))
	}
	return strings.Join(matches, "\n"), map[string]any{
		"path": relative, "pattern": args.Pattern, "include": args.Include,
		"matches": min(len(matches), s.Policy.MaxSearchMatches), "files_scanned": filesScanned,
		"files_skipped": filesSkipped, "match_limit_reached": limitReached,
	}, nil
}
