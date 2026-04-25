package types

import "github.com/google/uuid"

type FindingID string

func NewFindingID() FindingID { return FindingID(uuid.New().String()) }
func (x FindingID) String() string { return string(x) }

type IssueID string

func NewIssueID() IssueID { return IssueID(uuid.New().String()) }
func (x IssueID) String() string { return string(x) }

type ScanSessionID string

func NewScanSessionID() ScanSessionID { return ScanSessionID(uuid.New().String()) }
func (x ScanSessionID) String() string { return string(x) }

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityNote    Severity = "note"
)

type Verdict string

const (
	VerdictTruePositive  Verdict = "true_positive"
	VerdictFalsePositive Verdict = "false_positive"
	VerdictUncertain     Verdict = "uncertain"
)
