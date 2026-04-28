package types

import "github.com/m-mizutani/goerr/v2"

var (
	ErrKeyRuleID        = goerr.NewTypedKey[string]("rule_id")
	ErrKeyFindingID     = goerr.NewTypedKey[FindingID]("finding_id")
	ErrKeyIssueID       = goerr.NewTypedKey[IssueID]("issue_id")
	ErrKeyScanSessionID = goerr.NewTypedKey[ScanSessionID]("scan_session_id")
	ErrKeyFilePath      = goerr.NewTypedKey[string]("file_path")
	ErrKeySeverity      = goerr.NewTypedKey[Severity]("severity")
	ErrKeyVerdict       = goerr.NewTypedKey[Verdict]("verdict")
)
