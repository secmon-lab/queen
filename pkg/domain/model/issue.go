package model

import (
	"time"

	"github.com/secmon-lab/queen/pkg/domain/types"
)

type Issue struct {
	ID        types.IssueID      `json:"id"`
	SessionID types.ScanSessionID `json:"session_id"`
	Findings  []*Finding         `json:"findings"`
	Triage    *TriageResult      `json:"triage,omitempty"`
}

type TriageResult struct {
	Verdict types.Verdict `json:"verdict"`
	Reason  string        `json:"reason"`
}

type ScanSession struct {
	ID          types.ScanSessionID `json:"session_id"`
	StartedAt   time.Time           `json:"started_at"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	SARIFPath   string              `json:"sarif_path"`
	RepoPath    string              `json:"repo_path"`
	Issues      []*Issue            `json:"issues"`
}
