package usecase

import (
	"github.com/m-mizutani/gollem"
	"github.com/secmon-lab/queen/pkg/domain/interfaces"
)

type UseCases struct {
	llmClient  gollem.LLMClient
	repository interfaces.Repository
	tools      []gollem.Tool
}

type Option func(*UseCases)

func WithLLMClient(c gollem.LLMClient) Option {
	return func(u *UseCases) { u.llmClient = c }
}

func WithRepository(r interfaces.Repository) Option {
	return func(u *UseCases) { u.repository = r }
}

func WithTools(tools ...gollem.Tool) Option {
	return func(u *UseCases) { u.tools = tools }
}

func New(opts ...Option) *UseCases {
	u := &UseCases{}
	for _, opt := range opts {
		opt(u)
	}
	return u
}
