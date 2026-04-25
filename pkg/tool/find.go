package tool

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/gollem"
)

type Find struct {
	repoPath string
}

func NewFind(repoPath string) *Find {
	return &Find{repoPath: repoPath}
}

func (t *Find) Spec() gollem.ToolSpec {
	return gollem.ToolSpec{
		Name:        "find",
		Description: "Search for files in the repository by glob pattern. Returns matching file paths relative to the repository root.",
		Parameters: map[string]*gollem.Parameter{
			"pattern": {
				Type:        gollem.TypeString,
				Description: "Glob pattern to match files (e.g., \"*.go\", \"pkg/**/*.go\", \"**/handler.go\")",
				Required:    true,
			},
		},
	}
}

func (t *Find) Run(ctx context.Context, args map[string]any) (map[string]any, error) {
	pattern, _ := args["pattern"].(string)
	if pattern == "" {
		return nil, fmt.Errorf("pattern is required")
	}

	repoAbs, err := filepath.Abs(t.repoPath)
	if err != nil {
		return nil, fmt.Errorf("invalid repo path: %w", err)
	}

	var matches []string
	err = filepath.Walk(repoAbs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(repoAbs, path)
		if err != nil {
			return nil
		}

		matched, err := filepath.Match(pattern, rel)
		if err == nil && matched {
			matches = append(matches, rel)
			return nil
		}

		matched, err = filepath.Match(pattern, filepath.Base(rel))
		if err == nil && matched {
			matches = append(matches, rel)
			return nil
		}

		if strings.Contains(pattern, "**") {
			globPattern := strings.ReplaceAll(pattern, "**", "*")
			parts := strings.Split(rel, string(filepath.Separator))
			for i := range parts {
				subPath := filepath.Join(parts[i:]...)
				matched, err = filepath.Match(globPattern, subPath)
				if err == nil && matched {
					matches = append(matches, rel)
					return nil
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return map[string]any{
		"files": strings.Join(matches, "\n"),
		"count": len(matches),
	}, nil
}
