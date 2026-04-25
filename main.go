package main

import (
	"context"
	"os"

	"github.com/secmon-lab/queen/pkg/cli"
	"github.com/secmon-lab/queen/pkg/utils/logging"
)

func main() {
	if err := cli.Run(context.Background(), os.Args); err != nil {
		logging.Default().Error("fatal", logging.ErrAttr(err))
		os.Exit(1)
	}
}
