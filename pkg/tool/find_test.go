package tool_test

import (
	"context"
	"strings"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/tool"
)

func TestFindGlob(t *testing.T) {
	f := tool.NewFind(testRepoPath)
	result := gt.R1(f.Run(context.Background(), map[string]any{
		"pattern": "*.go",
	})).NoError(t)

	files := result["files"].(string)
	count := result["count"].(int)
	gt.N(t, count).GreaterOrEqual(2)
	gt.B(t, strings.Contains(files, "query.go")).True()
	gt.B(t, strings.Contains(files, "handler.go")).True()
}

func TestFindNestedGlob(t *testing.T) {
	f := tool.NewFind(testRepoPath)
	result := gt.R1(f.Run(context.Background(), map[string]any{
		"pattern": "pkg/**/*.go",
	})).NoError(t)

	count := result["count"].(int)
	gt.N(t, count).GreaterOrEqual(2)
}

func TestFindNoMatch(t *testing.T) {
	f := tool.NewFind(testRepoPath)
	result := gt.R1(f.Run(context.Background(), map[string]any{
		"pattern": "*.xyz",
	})).NoError(t)

	count := result["count"].(int)
	gt.Equal(t, count, 0)
}
