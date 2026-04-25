package fs_test

import (
	"context"
	"testing"
	"time"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
	"github.com/secmon-lab/queen/pkg/repository/fs"
)

func TestScanSessionRoundTrip(t *testing.T) {
	repo := fs.New(t.TempDir())
	ctx := context.Background()

	now := time.Now().Truncate(time.Second)
	session := &model.ScanSession{
		ID:        types.NewScanSessionID(),
		StartedAt: now,
		SARIFPath: "/path/to/sarif.json",
		RepoPath:  "/path/to/repo",
	}

	gt.NoError(t, repo.PutScanSession(ctx, session))

	got := gt.R1(repo.GetScanSession(ctx, session.ID)).NoError(t)
	gt.Equal(t, got.ID, session.ID)
	gt.Equal(t, got.SARIFPath, session.SARIFPath)
	gt.Equal(t, got.RepoPath, session.RepoPath)
	gt.Equal(t, got.StartedAt.Unix(), session.StartedAt.Unix())
}

func TestIssueRoundTrip(t *testing.T) {
	repo := fs.New(t.TempDir())
	ctx := context.Background()

	sessionID := types.NewScanSessionID()
	issue := &model.Issue{
		ID:        types.NewIssueID(),
		SessionID: sessionID,
		Findings: []*model.Finding{
			{
				ID:       types.NewFindingID(),
				RuleID:   "go.lang.security.audit.sqli",
				Severity: types.SeverityError,
				Message:  "SQL injection",
				Location: model.Location{
					FilePath:  "pkg/db/query.go",
					StartLine: 15,
					EndLine:   15,
				},
			},
		},
		Triage: &model.TriageResult{
			Verdict: types.VerdictTruePositive,
			Reason:  "User input directly concatenated",
		},
	}

	gt.NoError(t, repo.PutIssue(ctx, issue))

	got := gt.R1(repo.GetIssue(ctx, issue.ID)).NoError(t)
	gt.Equal(t, got.ID, issue.ID)
	gt.Equal(t, got.SessionID, sessionID)
	gt.Equal(t, len(got.Findings), 1)
	gt.Equal(t, got.Findings[0].RuleID, "go.lang.security.audit.sqli")
	gt.Equal(t, got.Triage.Verdict, types.VerdictTruePositive)
	gt.Equal(t, got.Triage.Reason, "User input directly concatenated")
}

func TestListIssuesSessionIsolation(t *testing.T) {
	repo := fs.New(t.TempDir())
	ctx := context.Background()

	session1 := types.NewScanSessionID()
	session2 := types.NewScanSessionID()

	for _, sid := range []types.ScanSessionID{session1, session1, session2} {
		issue := &model.Issue{
			ID:        types.NewIssueID(),
			SessionID: sid,
			Findings:  []*model.Finding{{ID: types.NewFindingID()}},
		}
		gt.NoError(t, repo.PutIssue(ctx, issue))
	}

	issues1 := gt.R1(repo.ListIssues(ctx, session1)).NoError(t)
	gt.Equal(t, len(issues1), 2)

	issues2 := gt.R1(repo.ListIssues(ctx, session2)).NoError(t)
	gt.Equal(t, len(issues2), 1)
}

func TestGetScanSessionNotFound(t *testing.T) {
	repo := fs.New(t.TempDir())
	ctx := context.Background()

	_, err := repo.GetScanSession(ctx, "nonexistent-id")
	gt.Error(t, err)
}

func TestGetIssueNotFound(t *testing.T) {
	repo := fs.New(t.TempDir())
	ctx := context.Background()

	_, err := repo.GetIssue(ctx, "nonexistent-id")
	gt.Error(t, err)
}
