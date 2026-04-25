package tool

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/gollem"
)

const maxFindResults = 100

var skipDirs = map[string]bool{
	".git":         true,
	"node_modules": true,
	"vendor":       true,
}

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
	err = filepath.WalkDir(repoAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		rel, err := filepath.Rel(repoAbs, path)
		if err != nil {
			return nil
		}

		if matchGlob(pattern, rel) {
			matches = append(matches, rel)
			if len(matches) >= maxFindResults {
				return fmt.Errorf("result limit reached")
			}
		}

		return nil
	})
	if err != nil && err.Error() != "result limit reached" {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	truncated := len(matches) >= maxFindResults

	return map[string]any{
		"files":     strings.Join(matches, "\n"),
		"count":     len(matches),
		"truncated": truncated,
	}, nil
}

func matchGlob(pattern, rel string) bool {
	if matched, err := filepath.Match(pattern, rel); err == nil && matched {
		return true
	}

	if matched, err := filepath.Match(pattern, filepath.Base(rel)); err == nil && matched {
		return true
	}

	if strings.Contains(pattern, "**") {
		return matchDoublestar(pattern, rel)
	}

	return false
}

func matchDoublestar(pattern, rel string) bool {
	parts := strings.SplitN(pattern, "**", 2)
	prefix := parts[0]
	suffix := parts[1]

	if prefix != "" {
		prefix = strings.TrimSuffix(prefix, string(filepath.Separator))
		if !strings.HasPrefix(rel, prefix+string(filepath.Separator)) && rel != prefix {
			return false
		}
		rel = strings.TrimPrefix(rel, prefix+string(filepath.Separator))
	}

	if suffix == "" {
		return true
	}
	suffix = strings.TrimPrefix(suffix, string(filepath.Separator))

	segments := strings.Split(rel, string(filepath.Separator))
	for i := range segments {
		candidate := filepath.Join(segments[i:]...)
		if matched, err := filepath.Match(suffix, candidate); err == nil && matched {
			return true
		}
	}

	if matched, err := filepath.Match(suffix, filepath.Base(rel)); err == nil && matched {
		return true
	}

	return false
}
