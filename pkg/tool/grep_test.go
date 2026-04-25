package tool_test

import (
	"context"
	"strings"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/tool"
)

func TestGrepMatch(t *testing.T) {
	g := tool.NewGrep(testRepoPath)
	result := gt.R1(g.Run(context.Background(), map[string]any{
		"pattern": "Sprintf",
	})).NoError(t)

	matches := result["matches"].(string)
	count := result["count"].(int)
	gt.N(t, count).GreaterOrEqual(2)
	gt.B(t, strings.Contains(matches, "pkg/db/query.go")).True()
}

func TestGrepWithPath(t *testing.T) {
	g := tool.NewGrep(testRepoPath)
	result := gt.R1(g.Run(context.Background(), map[string]any{
		"pattern": "exec.Command",
		"path":    "pkg/api",
	})).NoError(t)

	matches := result["matches"].(string)
	count := result["count"].(int)
	gt.N(t, count).GreaterOrEqual(2)
	gt.B(t, strings.Contains(matches, "handler.go")).True()
}

func TestGrepNoMatch(t *testing.T) {
	g := tool.NewGrep(testRepoPath)
	result := gt.R1(g.Run(context.Background(), map[string]any{
		"pattern": "xyznonexistent123",
	})).NoError(t)

	count := result["count"].(int)
	gt.Equal(t, count, 0)
}

func TestGrepMultipleFiles(t *testing.T) {
	g := tool.NewGrep(testRepoPath)
	result := gt.R1(g.Run(context.Background(), map[string]any{
		"pattern": "http.Request",
	})).NoError(t)

	matches := result["matches"].(string)
	gt.B(t, strings.Contains(matches, "pkg/db/query.go")).True()
	gt.B(t, strings.Contains(matches, "pkg/api/handler.go")).True()
}
