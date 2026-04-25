package tool

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/gollem"
)

type ReadFile struct {
	repoPath string
}

func NewReadFile(repoPath string) *ReadFile {
	return &ReadFile{repoPath: repoPath}
}

func (t *ReadFile) Spec() gollem.ToolSpec {
	return gollem.ToolSpec{
		Name:        "read_file",
		Description: "Read file contents from the repository. Returns lines with line numbers. Optionally specify a line range.",
		Parameters: map[string]*gollem.Parameter{
			"path": {
				Type:        gollem.TypeString,
				Description: "File path relative to the repository root",
				Required:    true,
			},
			"start_line": {
				Type:        gollem.TypeInteger,
				Description: "Start line number (1-based, inclusive). Omit to read from the beginning.",
			},
			"end_line": {
				Type:        gollem.TypeInteger,
				Description: "End line number (1-based, inclusive). Omit to read to the end.",
			},
		},
	}
}

func (t *ReadFile) Run(ctx context.Context, args map[string]any) (map[string]any, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	fullPath, err := t.safePath(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", path, err)
	}
	defer f.Close()

	startLine := 0
	endLine := 0
	if v, ok := toInt(args["start_line"]); ok {
		startLine = v
	}
	if v, ok := toInt(args["end_line"]); ok {
		endLine = v
	}

	var lines []string
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if startLine > 0 && lineNum < startLine {
			continue
		}
		if endLine > 0 && lineNum > endLine {
			break
		}
		lines = append(lines, fmt.Sprintf("%d: %s", lineNum, scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	return map[string]any{
		"content": strings.Join(lines, "\n"),
	}, nil
}

func (t *ReadFile) safePath(path string) (string, error) {
	joined := filepath.Join(t.repoPath, path)
	abs, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("invalid path: %w", err)
	}
	repoAbs, err := filepath.Abs(t.repoPath)
	if err != nil {
		return "", fmt.Errorf("invalid repo path: %w", err)
	}
	if !strings.HasPrefix(abs, repoAbs+string(filepath.Separator)) && abs != repoAbs {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}
	return abs, nil
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case float64:
		return int(n), true
	case int:
		return n, true
	case int64:
		return int(n), true
	default:
		return 0, false
	}
}
