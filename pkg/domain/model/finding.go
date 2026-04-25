package model

import "github.com/secmon-lab/queen/pkg/domain/types"

type Finding struct {
	ID       types.FindingID `json:"id"`
	RuleID   string          `json:"rule_id"`
	Severity types.Severity  `json:"severity"`
	Message  string          `json:"message"`
	Location Location        `json:"location"`
}

type Location struct {
	FilePath  string `json:"file_path"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}
