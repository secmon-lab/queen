package tool_test

import (
	"context"
	"strings"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/tool"
)

const testRepoPath = "../../examples/sample-repo"

func TestReadFileWholeFile(t *testing.T) {
	rf := tool.NewReadFile(testRepoPath)
	result := gt.R1(rf.Run(context.Background(), map[string]any{
		"path": "pkg/db/query.go",
	})).NoError(t)

	content := result["content"].(string)
	gt.B(t, strings.Contains(content, "package db")).True()
	gt.B(t, strings.Contains(content, "1: ")).True()
}

func TestReadFileLineRange(t *testing.T) {
	rf := tool.NewReadFile(testRepoPath)
	result := gt.R1(rf.Run(context.Background(), map[string]any{
		"path":       "pkg/db/query.go",
		"start_line": float64(16),
		"end_line":   float64(19),
	})).NoError(t)

	content := result["content"].(string)
	lines := strings.Split(content, "\n")
	gt.Equal(t, len(lines), 4)
	gt.B(t, strings.HasPrefix(lines[0], "16: ")).True()
	gt.B(t, strings.HasPrefix(lines[3], "19: ")).True()
}

func TestReadFileBoundary(t *testing.T) {
	rf := tool.NewReadFile(testRepoPath)
	result := gt.R1(rf.Run(context.Background(), map[string]any{
		"path":       "pkg/db/query.go",
		"start_line": float64(1),
		"end_line":   float64(2),
	})).NoError(t)

	content := result["content"].(string)
	lines := strings.Split(content, "\n")
	gt.Equal(t, len(lines), 2)
	gt.B(t, strings.HasPrefix(lines[0], "1: ")).True()
}

func TestReadFileNotFound(t *testing.T) {
	rf := tool.NewReadFile(testRepoPath)
	_, err := rf.Run(context.Background(), map[string]any{
		"path": "nonexistent.go",
	})
	gt.Error(t, err)
}

func TestReadFilePathTraversal(t *testing.T) {
	rf := tool.NewReadFile(testRepoPath)
	_, err := rf.Run(context.Background(), map[string]any{
		"path": "../../../etc/passwd",
	})
	gt.Error(t, err)
}
