package interfaces

import (
	"context"

	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
)

type Repository interface {
	PutScanSession(ctx context.Context, session *model.ScanSession) error
	GetScanSession(ctx context.Context, id types.ScanSessionID) (*model.ScanSession, error)

	PutIssue(ctx context.Context, issue *model.Issue) error
	GetIssue(ctx context.Context, id types.IssueID) (*model.Issue, error)
	ListIssues(ctx context.Context, sessionID types.ScanSessionID) ([]*model.Issue, error)
}
