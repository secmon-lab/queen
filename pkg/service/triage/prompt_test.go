package triage_test

import (
	"strings"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
	"github.com/secmon-lab/queen/pkg/service/triage"
)

func TestSystemPrompt(t *testing.T) {
	prompt := triage.SystemPrompt()
	gt.B(t, strings.Contains(prompt, "true_positive")).True()
	gt.B(t, strings.Contains(prompt, "false_positive")).True()
	gt.B(t, strings.Contains(prompt, "uncertain")).True()
	gt.B(t, strings.Contains(prompt, "security")).True()
}

func TestBuildPrompt(t *testing.T) {
	issue := &model.Issue{
		ID: types.NewIssueID(),
		Findings: []*model.Finding{
			{
				ID:       types.NewFindingID(),
				RuleID:   "go.lang.security.audit.sqli",
				Severity: types.SeverityError,
				Message:  "Possible SQL injection",
				Location: model.Location{
					FilePath:  "pkg/db/query.go",
					StartLine: 17,
					EndLine:   17,
				},
			},
		},
	}

	prompt := triage.BuildPrompt(issue)
	gt.B(t, strings.Contains(prompt, "go.lang.security.audit.sqli")).True()
	gt.B(t, strings.Contains(prompt, "error")).True()
	gt.B(t, strings.Contains(prompt, "Possible SQL injection")).True()
	gt.B(t, strings.Contains(prompt, "pkg/db/query.go")).True()
	gt.B(t, strings.Contains(prompt, "17")).True()
}

func TestBuildPromptMultipleFindings(t *testing.T) {
	issue := &model.Issue{
		ID: types.NewIssueID(),
		Findings: []*model.Finding{
			{
				ID:       types.NewFindingID(),
				RuleID:   "rule-a",
				Severity: types.SeverityError,
				Message:  "First finding",
				Location: model.Location{FilePath: "a.go", StartLine: 1, EndLine: 1},
			},
			{
				ID:       types.NewFindingID(),
				RuleID:   "rule-b",
				Severity: types.SeverityWarning,
				Message:  "Second finding",
				Location: model.Location{FilePath: "b.go", StartLine: 10, EndLine: 10},
			},
		},
	}

	prompt := triage.BuildPrompt(issue)
	gt.B(t, strings.Contains(prompt, "rule-a")).True()
	gt.B(t, strings.Contains(prompt, "rule-b")).True()
	gt.B(t, strings.Contains(prompt, "First finding")).True()
	gt.B(t, strings.Contains(prompt, "Second finding")).True()
	gt.B(t, strings.Contains(prompt, "Finding 1")).True()
	gt.B(t, strings.Contains(prompt, "Finding 2")).True()
}
