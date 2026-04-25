package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/secmon-lab/queen/pkg/cli"
)

func main() {
	if err := cli.Run(context.Background(), os.Args); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}
