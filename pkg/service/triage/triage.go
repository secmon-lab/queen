package triage

import (
	"context"
	"encoding/json"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/gollem"
	"github.com/secmon-lab/queen/pkg/domain/model"
	"github.com/secmon-lab/queen/pkg/domain/types"
)

type triageResponse struct {
	Verdict string `json:"verdict" description:"One of: true_positive, false_positive, uncertain"`
	Reason  string `json:"reason" description:"Concise explanation for the verdict"`
}

type Service struct {
	llmClient gollem.LLMClient
	tools     []gollem.Tool
}

func New(llmClient gollem.LLMClient, tools ...gollem.Tool) *Service {
	return &Service{
		llmClient: llmClient,
		tools:     tools,
	}
}

func (s *Service) Triage(ctx context.Context, issue *model.Issue) (*model.TriageResult, error) {
	userPrompt := BuildPrompt(issue)

	agent := gollem.New(
		s.llmClient,
		gollem.WithSystemPrompt(SystemPrompt()),
		gollem.WithTools(s.tools...),
	)

	resp, err := agent.Execute(ctx, gollem.Text(userPrompt))
	if err != nil {
		return nil, goerr.Wrap(err, "agent execution failed")
	}

	// Try to parse the agent's final response as structured JSON
	for _, text := range resp.Texts {
		var tr triageResponse
		if err := json.Unmarshal([]byte(text), &tr); err == nil && tr.Verdict != "" {
			return toTriageResult(&tr), nil
		}
	}

	// Fallback: use Query for structured output
	queryResp, err := gollem.Query[triageResponse](
		ctx,
		s.llmClient,
		"Based on the analysis above, provide your final verdict.\n\nContext from previous analysis:\n"+resp.String(),
		gollem.WithQuerySystemPrompt(SystemPrompt()),
	)
	if err != nil {
		return nil, goerr.Wrap(err, "failed to get structured triage result")
	}

	return toTriageResult(queryResp.Data), nil
}

func toTriageResult(tr *triageResponse) *model.TriageResult {
	var verdict types.Verdict
	switch tr.Verdict {
	case "true_positive":
		verdict = types.VerdictTruePositive
	case "false_positive":
		verdict = types.VerdictFalsePositive
	default:
		verdict = types.VerdictUncertain
	}

	return &model.TriageResult{
		Verdict: verdict,
		Reason:  tr.Reason,
	}
}
