package sarif_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/domain/types"
	"github.com/secmon-lab/queen/pkg/service/sarif"
)

func TestParse(t *testing.T) {
	findings := gt.R1(sarif.Parse("../../../examples/sarif/semgrep-result.sarif")).NoError(t)
	gt.Equal(t, len(findings), 5)

	gt.Equal(t, findings[0].RuleID, "go.lang.security.audit.sqli.taint-sql-string-format")
	gt.Equal(t, findings[0].Severity, types.SeverityError)
	gt.Equal(t, findings[0].Location.FilePath, "pkg/db/query.go")
	gt.Equal(t, findings[0].Location.StartLine, 17)

	gt.Equal(t, findings[1].RuleID, "go.lang.security.audit.sqli.taint-sql-string-format")
	gt.Equal(t, findings[1].Location.StartLine, 25)

	gt.Equal(t, findings[2].RuleID, "go.lang.security.audit.command-injection")
	gt.Equal(t, findings[2].Severity, types.SeverityError)
	gt.Equal(t, findings[2].Location.FilePath, "pkg/api/handler.go")
	gt.Equal(t, findings[2].Location.StartLine, 12)

	gt.Equal(t, findings[3].RuleID, "go.lang.security.audit.command-injection")
	gt.Equal(t, findings[3].Severity, types.SeverityWarning)
	gt.Equal(t, findings[3].Location.StartLine, 22)

	gt.Equal(t, findings[4].RuleID, "go.lang.security.audit.xss.direct-response-write")
	gt.Equal(t, findings[4].Severity, types.SeverityError)
	gt.Equal(t, findings[4].Location.StartLine, 33)
}

func TestParseInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "bad.sarif")
	gt.NoError(t, os.WriteFile(path, []byte("{invalid json"), 0o644))

	_, err := sarif.Parse(path)
	gt.Error(t, err)
}

func TestParseEmptyResults(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "empty.sarif")
	data := `{"version":"2.1.0","runs":[{"tool":{"driver":{"name":"test"}},"results":[]}]}`
	gt.NoError(t, os.WriteFile(path, []byte(data), 0o644))

	findings := gt.R1(sarif.Parse(path)).NoError(t)
	gt.Equal(t, len(findings), 0)
}

func TestParseFileNotFound(t *testing.T) {
	_, err := sarif.Parse("/nonexistent/path.sarif")
	gt.Error(t, err)
}
