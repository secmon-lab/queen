package config

import (
	"context"

	"github.com/m-mizutani/goerr/v2"
	"github.com/m-mizutani/gollem"
	"github.com/m-mizutani/gollem/llm/claude"
	"github.com/m-mizutani/gollem/llm/gemini"
	"github.com/m-mizutani/gollem/llm/openai"
	"github.com/urfave/cli/v3"
)

type LLMCfg struct {
	openaiAPIKey string
	openaiModel  string

	anthropicAPIKey string
	anthropicModel  string

	geminiProjectID string
	geminiLocation  string
	geminiModel     string
}

func (x *LLMCfg) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "openai-api-key",
			Usage:       "OpenAI API key",
			Sources:     cli.EnvVars("QUEEN_OPENAI_API_KEY"),
			Destination: &x.openaiAPIKey,
			Category:    "OpenAI",
		},
		&cli.StringFlag{
			Name:        "openai-model",
			Usage:       "OpenAI model name",
			Value:       "gpt-4o",
			Sources:     cli.EnvVars("QUEEN_OPENAI_MODEL"),
			Destination: &x.openaiModel,
			Category:    "OpenAI",
		},
		&cli.StringFlag{
			Name:        "anthropic-api-key",
			Usage:       "Anthropic API key",
			Sources:     cli.EnvVars("QUEEN_ANTHROPIC_API_KEY"),
			Destination: &x.anthropicAPIKey,
			Category:    "Anthropic",
		},
		&cli.StringFlag{
			Name:        "anthropic-model",
			Usage:       "Anthropic model name",
			Value:       "claude-sonnet-4-20250514",
			Sources:     cli.EnvVars("QUEEN_ANTHROPIC_MODEL"),
			Destination: &x.anthropicModel,
			Category:    "Anthropic",
		},
		&cli.StringFlag{
			Name:        "gemini-project-id",
			Usage:       "Google Cloud project ID for Gemini",
			Sources:     cli.EnvVars("QUEEN_GEMINI_PROJECT_ID"),
			Destination: &x.geminiProjectID,
			Category:    "Gemini",
		},
		&cli.StringFlag{
			Name:        "gemini-location",
			Usage:       "Google Cloud location for Gemini",
			Value:       "us-central1",
			Sources:     cli.EnvVars("QUEEN_GEMINI_LOCATION"),
			Destination: &x.geminiLocation,
			Category:    "Gemini",
		},
		&cli.StringFlag{
			Name:        "gemini-model",
			Usage:       "Gemini model name",
			Value:       "gemini-2.5-flash",
			Sources:     cli.EnvVars("QUEEN_GEMINI_MODEL"),
			Destination: &x.geminiModel,
			Category:    "Gemini",
		},
	}
}

func (x *LLMCfg) Configure(ctx context.Context) (gollem.LLMClient, error) {
	switch {
	case x.anthropicAPIKey != "":
		client, err := claude.New(ctx, x.anthropicAPIKey, claude.WithModel(x.anthropicModel))
		if err != nil {
			return nil, goerr.Wrap(err, "failed to create Anthropic client")
		}
		return client, nil

	case x.openaiAPIKey != "":
		client, err := openai.New(ctx, x.openaiAPIKey, openai.WithModel(x.openaiModel))
		if err != nil {
			return nil, goerr.Wrap(err, "failed to create OpenAI client")
		}
		return client, nil

	case x.geminiProjectID != "":
		client, err := gemini.New(ctx, x.geminiProjectID, x.geminiLocation, gemini.WithModel(x.geminiModel))
		if err != nil {
			return nil, goerr.Wrap(err, "failed to create Gemini client")
		}
		return client, nil

	default:
		return nil, goerr.New("no LLM provider configured: set one of --anthropic-api-key, --openai-api-key, or --gemini-project-id")
	}
}
