package cli

import (
	"context"

	"github.com/secmon-lab/queen/pkg/cli/config"
	"github.com/secmon-lab/queen/pkg/utils/logging"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, args []string) error {
	var logCfg config.Logger

	app := &cli.Command{
		Name:  "queen",
		Usage: "An agentic SAST triage tool",
		Flags: logCfg.Flags(),
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			logger, closer, err := logCfg.Configure()
			if err != nil {
				return ctx, err
			}
			_ = closer

			ctx = logging.With(ctx, logger)
			return ctx, nil
		},
		Commands: []*cli.Command{
			cmdScan(),
		},
	}

	return app.Run(ctx, args)
}
