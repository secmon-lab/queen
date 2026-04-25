package usecase

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
	"github.com/secmon-lab/queen/pkg/service/sarif"
	"github.com/secmon-lab/queen/pkg/service/triage"
)

func (u *UseCases) Scan(ctx context.Context, sarifPath, repoPath string) (*model.ScanSession, error) {
	if _, err := os.Stat(sarifPath); err != nil {
		return nil, goerr.Wrap(err, "SARIF file not accessible", goerr.V("path", sarifPath))
	}
	if _, err := os.Stat(repoPath); err != nil {
		return nil, goerr.Wrap(err, "repository path not accessible", goerr.V("path", repoPath))
	}

	session := &model.ScanSession{
		ID:        types.NewScanSessionID(),
		StartedAt: time.Now(),
		SARIFPath: sarifPath,
		RepoPath:  repoPath,
	}

	if u.repository != nil {
		if err := u.repository.PutScanSession(ctx, session); err != nil {
			return nil, goerr.Wrap(err, "failed to save scan session")
		}
	}

	findings, err := sarif.Parse(sarifPath)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to parse SARIF")
	}

	triageSvc := triage.New(u.llmClient, u.tools...)

	for _, finding := range findings {
		issue := &model.Issue{
			ID:        types.NewIssueID(),
			SessionID: session.ID,
			Findings:  []*model.Finding{finding},
		}

		result, err := triageSvc.Triage(ctx, issue)
		if err != nil {
			slog.ErrorContext(ctx, "triage failed, skipping issue",
				"issue_id", issue.ID,
				"rule_id", finding.RuleID,
				"error", err,
			)
		} else {
			issue.Triage = result
		}

		if u.repository != nil {
			if err := u.repository.PutIssue(ctx, issue); err != nil {
				slog.ErrorContext(ctx, "failed to save issue",
					"issue_id", issue.ID,
					"error", err,
				)
			}
		}

		session.Issues = append(session.Issues, issue)
	}

	now := time.Now()
	session.CompletedAt = &now

	if u.repository != nil {
		if err := u.repository.PutScanSession(ctx, session); err != nil {
			return nil, goerr.Wrap(err, "failed to update scan session")
		}
	}

	return session, nil
}
