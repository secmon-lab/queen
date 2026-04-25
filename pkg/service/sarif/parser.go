package sarif

import (
	"encoding/json"
	"os"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
)

func Parse(filePath string) ([]*model.Finding, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to read SARIF file", goerr.V("path", filePath))
	}

	var report model.SARIFReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, goerr.Wrap(err, "failed to parse SARIF file")
	}

	var findings []*model.Finding
	for _, run := range report.Runs {
		for _, result := range run.Results {
			finding := &model.Finding{
				ID:       types.NewFindingID(),
				RuleID:   result.RuleID,
				Severity: toSeverity(result.Level),
				Message:  result.Message.Text,
			}

			if len(result.Locations) > 0 {
				loc := result.Locations[0].PhysicalLocation
				endLine := loc.Region.EndLine
				if endLine == 0 {
					endLine = loc.Region.StartLine
				}
				finding.Location = model.Location{
					FilePath:  loc.ArtifactLocation.URI,
					StartLine: loc.Region.StartLine,
					EndLine:   endLine,
				}
			}

			findings = append(findings, finding)
		}
	}

	return findings, nil
}

func toSeverity(level string) types.Severity {
	switch level {
	case "error":
		return types.SeverityError
	case "warning":
		return types.SeverityWarning
	case "note":
		return types.SeverityNote
	default:
		return types.SeverityWarning
	}
}
