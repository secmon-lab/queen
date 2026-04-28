package logging_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/m-mizutani/gt"
	"github.com/secmon-lab/queen/pkg/utils/logging"
)

func TestNewJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(&buf, slog.LevelInfo, logging.FormatJSON, false)
	logger.Info("hello", slog.String("key", "value"))
	gt.S(t, buf.String()).Contains(`"message":"hello"`).Contains(`"key":"value"`)
}

func TestNewConsole(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(&buf, slog.LevelInfo, logging.FormatConsole, false)
	logger.Info("hello")
	gt.S(t, buf.String()).Contains("hello")
}

func TestSecretMasking(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(&buf, slog.LevelInfo, logging.FormatJSON, false)
	logger.Info("test",
		slog.String("secret_key", "mysecret"),
		slog.String("normal_key", "visible"),
	)
	gt.S(t, buf.String()).Contains("visible").NotContains("mysecret")
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(&buf, slog.LevelWarn, logging.FormatJSON, false)
	logger.Info("should not appear")
	gt.S(t, buf.String()).Equal("")

	logger.Warn("should appear")
	gt.S(t, buf.String()).Contains("should appear")
}

func TestContextPropagation(t *testing.T) {
	var buf bytes.Buffer
	logger := logging.New(&buf, slog.LevelInfo, logging.FormatJSON, false)
	ctx := logging.With(context.Background(), logger)

	retrieved := logging.From(ctx)
	retrieved.Info("from context")
	gt.S(t, buf.String()).Contains("from context")
}

func TestContextFallbackToDefault(t *testing.T) {
	logger := logging.From(context.Background())
	gt.V(t, logger).NotNil()
}

func TestSetDefault(t *testing.T) {
	original := logging.Default()
	defer logging.SetDefault(original)

	var buf bytes.Buffer
	logger := logging.New(&buf, slog.LevelInfo, logging.FormatJSON, false)
	logging.SetDefault(logger)

	gt.V(t, logging.Default()).Equal(logger)

	logging.From(context.Background()).Info("via default")
	gt.S(t, buf.String()).Contains("via default")
}
