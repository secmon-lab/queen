package triage

import (
	"fmt"
	"strings"

	"github.com/secmon-lab/queen/pkg/domain/model"
)

func SystemPrompt() string {
	return `You are an expert security engineer performing triage on Static Application Security Testing (SAST) findings.

Your task is to analyze each finding and determine whether it represents a real vulnerability (true positive) or a false alarm (false positive).

For each finding, you will be given:
- The SAST rule ID and severity
- The detection message from the SAST tool
- The file path and line numbers where the issue was detected

You have access to tools to read source code files, search for files, and grep for patterns. Use these tools to:
1. Read the flagged source code and its surrounding context
2. Trace data flow to understand if user-controlled input reaches the vulnerable sink
3. Check for existing sanitization, validation, or safe usage patterns
4. Look at related files if needed to understand the full picture

After your analysis, provide your verdict:
- "true_positive": The finding represents a real, exploitable vulnerability
- "false_positive": The finding is a false alarm (input is not user-controlled, properly sanitized, etc.)
- "uncertain": You cannot determine with confidence whether it's a real vulnerability

Provide a clear, concise reason for your verdict.`
}

func BuildPrompt(issue *model.Issue) string {
	var b strings.Builder

	fmt.Fprintf(&b, "Please triage the following SAST finding(s):\n\n")

	for i, f := range issue.Findings {
		if len(issue.Findings) > 1 {
			fmt.Fprintf(&b, "### Finding %d\n", i+1)
		}
		fmt.Fprintf(&b, "- **Rule ID**: %s\n", f.RuleID)
		fmt.Fprintf(&b, "- **Severity**: %s\n", f.Severity)
		fmt.Fprintf(&b, "- **Message**: %s\n", f.Message)
		fmt.Fprintf(&b, "- **File**: %s (lines %d-%d)\n", f.Location.FilePath, f.Location.StartLine, f.Location.EndLine)
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "Use the available tools to read the source code and investigate. Then provide your verdict and reason.")

	return b.String()
}
