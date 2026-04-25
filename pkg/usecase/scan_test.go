package usecase_test

import (
	"context"
	"testing"

	"github.com/m-mizutani/gollem"
	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/domain/types"
	"github.com/secmon-lab/queen/pkg/repository/fs"
	"github.com/secmon-lab/queen/pkg/usecase"
)

type mockSession struct{}

func (s *mockSession) Generate(ctx context.Context, input []gollem.Input, opts ...gollem.GenerateOption) (*gollem.Response, error) {
	return &gollem.Response{
		Texts: []string{`{"verdict":"true_positive","reason":"test reason"}`},
	}, nil
}

func (s *mockSession) Stream(ctx context.Context, input []gollem.Input, opts ...gollem.GenerateOption) (<-chan *gollem.Response, error) {
	ch := make(chan *gollem.Response, 1)
	ch <- &gollem.Response{Texts: []string{`{"verdict":"true_positive","reason":"test"}`}}
	close(ch)
	return ch, nil
}

func (s *mockSession) GenerateContent(ctx context.Context, input ...gollem.Input) (*gollem.Response, error) {
	return s.Generate(ctx, input)
}

func (s *mockSession) GenerateStream(ctx context.Context, input ...gollem.Input) (<-chan *gollem.Response, error) {
	return s.Stream(ctx, input)
}

func (s *mockSession) History() (*gollem.History, error) {
	return &gollem.History{}, nil
}

func (s *mockSession) AppendHistory(h *gollem.History) error {
	return nil
}

func (s *mockSession) CountToken(ctx context.Context, input ...gollem.Input) (int, error) {
	return 0, nil
}

type mockLLMClient struct{}

func (c *mockLLMClient) NewSession(ctx context.Context, opts ...gollem.SessionOption) (gollem.Session, error) {
	return &mockSession{}, nil
}

func (c *mockLLMClient) GenerateEmbedding(ctx context.Context, dimension int, input []string) ([][]float64, error) {
	return nil, nil
}

func TestScanFullFlow(t *testing.T) {
	repo := fs.New(t.TempDir())
	uc := usecase.New(
		usecase.WithLLMClient(&mockLLMClient{}),
		usecase.WithRepository(repo),
	)

	session := gt.R1(uc.Scan(context.Background(),
		"../../testdata/sarif/semgrep-result.sarif",
		"../../testdata/sample-repo",
	)).NoError(t)

	gt.Equal(t, len(session.Issues), 5)
	for _, issue := range session.Issues {
		gt.Equal(t, len(issue.Findings), 1)
		gt.B(t, issue.Triage != nil).True()
		gt.Equal(t, issue.Triage.Verdict, types.VerdictTruePositive)
	}
	gt.B(t, session.CompletedAt != nil).True()
}

func TestScanSARIFNotFound(t *testing.T) {
	uc := usecase.New(
		usecase.WithLLMClient(&mockLLMClient{}),
	)

	_, err := uc.Scan(context.Background(), "/nonexistent.sarif", "../../testdata/sample-repo")
	gt.Error(t, err)
}

func TestScanRepoNotFound(t *testing.T) {
	uc := usecase.New(
		usecase.WithLLMClient(&mockLLMClient{}),
	)

	_, err := uc.Scan(context.Background(), "../../testdata/sarif/semgrep-result.sarif", "/nonexistent/repo")
	gt.Error(t, err)
}

func TestScanLLMError(t *testing.T) {
	repo := fs.New(t.TempDir())
	uc := usecase.New(
		usecase.WithLLMClient(&errorLLMClient{}),
		usecase.WithRepository(repo),
	)

	session := gt.R1(uc.Scan(context.Background(),
		"../../testdata/sarif/semgrep-result.sarif",
		"../../testdata/sample-repo",
	)).NoError(t)

	gt.Equal(t, len(session.Issues), 5)
	for _, issue := range session.Issues {
		gt.B(t, issue.Triage == nil).True()
	}
}

type errorSession struct{}

func (s *errorSession) Generate(ctx context.Context, input []gollem.Input, opts ...gollem.GenerateOption) (*gollem.Response, error) {
	return nil, gollem.ErrExitConversation
}

func (s *errorSession) Stream(ctx context.Context, input []gollem.Input, opts ...gollem.GenerateOption) (<-chan *gollem.Response, error) {
	return nil, gollem.ErrExitConversation
}

func (s *errorSession) GenerateContent(ctx context.Context, input ...gollem.Input) (*gollem.Response, error) {
	return nil, gollem.ErrExitConversation
}

func (s *errorSession) GenerateStream(ctx context.Context, input ...gollem.Input) (<-chan *gollem.Response, error) {
	return nil, gollem.ErrExitConversation
}

func (s *errorSession) History() (*gollem.History, error) {
	return &gollem.History{}, nil
}

func (s *errorSession) AppendHistory(h *gollem.History) error {
	return nil
}

func (s *errorSession) CountToken(ctx context.Context, input ...gollem.Input) (int, error) {
	return 0, nil
}

type errorLLMClient struct{}

func (c *errorLLMClient) NewSession(ctx context.Context, opts ...gollem.SessionOption) (gollem.Session, error) {
	return &errorSession{}, nil
}

func (c *errorLLMClient) GenerateEmbedding(ctx context.Context, dimension int, input []string) ([][]float64, error) {
	return nil, nil
}
