package cli

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, args []string) error {
	var logLevel string

	app := &cli.Command{
		Name:  "queen",
		Usage: "An agentic SAST triage tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "Log level (debug, info, warn, error)",
				Value:       "info",
				Sources:     cli.EnvVars("QUEEN_LOG_LEVEL"),
				Destination: &logLevel,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			var level slog.Level
			switch logLevel {
			case "debug":
				level = slog.LevelDebug
			case "warn":
				level = slog.LevelWarn
			case "error":
				level = slog.LevelError
			default:
				level = slog.LevelInfo
			}
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: level})))
			return ctx, nil
		},
		Commands: []*cli.Command{
			cmdScan(),
		},
	}

	return app.Run(ctx, args)
}
