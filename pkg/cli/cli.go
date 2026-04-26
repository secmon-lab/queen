package cli

import (
	"context"

	"github.com/secmon-lab/queen/pkg/cli/config"
	"github.com/secmon-lab/queen/pkg/utils/logging"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, args []string) error {
	var logCfg config.Logger
	closer := func() {}

	app := &cli.Command{
		Name:  "queen",
		Usage: "An agentic SAST triage tool",
		Flags: logCfg.Flags(),
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			logger, c2, err := logCfg.Configure()
			if err != nil {
				return ctx, err
			}
			closer = c2

			ctx = logging.With(ctx, logger)
			return ctx, nil
		},
		After: func(ctx context.Context, c *cli.Command) error {
			closer()
			return nil
		},
		Commands: []*cli.Command{
			cmdScan(),
		},
	}

	return app.Run(ctx, args)
}
