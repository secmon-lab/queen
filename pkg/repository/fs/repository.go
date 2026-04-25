package fs

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/queen/pkg/domain/interfaces"
	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
)

var _ interfaces.Repository = &Repository{}

type Repository struct {
	baseDir string
}

func New(baseDir string) *Repository {
	return &Repository{baseDir: baseDir}
}

func (r *Repository) PutScanSession(ctx context.Context, session *model.ScanSession) error {
	dir := filepath.Join(r.baseDir, "sessions")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return goerr.Wrap(err, "failed to create sessions directory")
	}

	data, err := json.Marshal(session)
	if err != nil {
		return goerr.Wrap(err, "failed to marshal scan session")
	}

	path := filepath.Join(dir, session.ID.String()+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return goerr.Wrap(err, "failed to write scan session")
	}

	return nil
}

func (r *Repository) GetScanSession(ctx context.Context, id types.ScanSessionID) (*model.ScanSession, error) {
	path := filepath.Join(r.baseDir, "sessions", id.String()+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, goerr.Wrap(err, "scan session not found", goerr.V("id", id))
		}
		return nil, goerr.Wrap(err, "failed to read scan session")
	}

	var session model.ScanSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, goerr.Wrap(err, "failed to unmarshal scan session")
	}

	return &session, nil
}

func (r *Repository) PutIssue(ctx context.Context, issue *model.Issue) error {
	dir := filepath.Join(r.baseDir, "issues", issue.SessionID.String())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return goerr.Wrap(err, "failed to create issues directory")
	}

	data, err := json.Marshal(issue)
	if err != nil {
		return goerr.Wrap(err, "failed to marshal issue")
	}

	path := filepath.Join(dir, issue.ID.String()+".json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return goerr.Wrap(err, "failed to write issue")
	}

	return nil
}

func (r *Repository) GetIssue(ctx context.Context, id types.IssueID) (*model.Issue, error) {
	pattern := filepath.Join(r.baseDir, "issues", "*", id.String()+".json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to glob issue files")
	}
	if len(matches) == 0 {
		return nil, goerr.New("issue not found", goerr.V("id", id))
	}

	data, err := os.ReadFile(matches[0])
	if err != nil {
		return nil, goerr.Wrap(err, "failed to read issue")
	}

	var issue model.Issue
	if err := json.Unmarshal(data, &issue); err != nil {
		return nil, goerr.Wrap(err, "failed to unmarshal issue")
	}

	return &issue, nil
}

func (r *Repository) ListIssues(ctx context.Context, sessionID types.ScanSessionID) ([]*model.Issue, error) {
	dir := filepath.Join(r.baseDir, "issues", sessionID.String())
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, goerr.Wrap(err, "failed to read issues directory")
	}

	var issues []*model.Issue
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, goerr.Wrap(err, "failed to read issue file")
		}

		var issue model.Issue
		if err := json.Unmarshal(data, &issue); err != nil {
			return nil, goerr.Wrap(err, "failed to unmarshal issue")
		}
		issues = append(issues, &issue)
	}

	return issues, nil
}
