package tool

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/gollem"
)

type Grep struct {
	repoPath string
}

func NewGrep(repoPath string) *Grep {
	return &Grep{repoPath: repoPath}
}

func (t *Grep) Spec() gollem.ToolSpec {
	return gollem.ToolSpec{
		Name:        "grep",
		Description: "Search for a keyword in files within the repository. Returns matching lines with file path and line number.",
		Parameters: map[string]*gollem.Parameter{
			"pattern": {
				Type:        gollem.TypeString,
				Description: "Search keyword or pattern (case-sensitive substring match)",
				Required:    true,
			},
			"path": {
				Type:        gollem.TypeString,
				Description: "Directory or file path to search in (relative to repo root). Omit to search entire repo.",
			},
		},
	}
}

func (t *Grep) Run(ctx context.Context, args map[string]any) (map[string]any, error) {
	pattern, _ := args["pattern"].(string)
	if pattern == "" {
		return nil, goerr.New("pattern is required")
	}

	repoAbs, err := filepath.Abs(t.repoPath)
	if err != nil {
		return nil, goerr.Wrap(err, "invalid repo path")
	}

	searchPath := repoAbs
	if p, ok := args["path"].(string); ok && p != "" {
		joined := filepath.Join(t.repoPath, p)
		abs, err := filepath.Abs(joined)
		if err != nil {
			return nil, goerr.Wrap(err, "invalid path")
		}
		if !strings.HasPrefix(abs, repoAbs+string(filepath.Separator)) && abs != repoAbs {
			return nil, goerr.New("path traversal detected", goerr.V("path", p))
		}
		searchPath = abs
	}

	var results []string
	err = filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		if info.Size() > 1<<20 {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		rel, err := filepath.Rel(repoAbs, path)
		if err != nil {
			return nil
		}

		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if strings.Contains(line, pattern) {
				results = append(results, fmt.Sprintf("%s:%d: %s", rel, lineNum, line))
			}
		}

		return nil
	})
	if err != nil {
		return nil, goerr.Wrap(err, "failed to walk directory")
	}

	return map[string]any{
		"matches": strings.Join(results, "\n"),
		"count":   len(results),
	}, nil
}
