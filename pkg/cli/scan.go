package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/queen/pkg/cli/config"
	"github.com/secmon-lab/queen/pkg/repository/fs"
	"github.com/secmon-lab/queen/pkg/tool"
	"github.com/secmon-lab/queen/pkg/usecase"
	"github.com/urfave/cli/v3"
)

func cmdScan() *cli.Command {
	var (
		llmCfg  config.LLMCfg
		sarif   string
		repo    string
		dataDir string
	)

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "sarif",
			Usage:       "Path to SARIF file",
			Required:    true,
			Destination: &sarif,
		},
		&cli.StringFlag{
			Name:        "repo",
			Usage:       "Path to target repository",
			Required:    true,
			Destination: &repo,
		},
		&cli.StringFlag{
			Name:        "data-dir",
			Usage:       "Directory for storing scan data",
			Value:       filepath.Join(os.Getenv("HOME"), ".queen", "data"),
			Sources:     cli.EnvVars("QUEEN_DATA_DIR"),
			Destination: &dataDir,
		},
	}
	flags = append(flags, llmCfg.Flags()...)

	return &cli.Command{
		Name:  "scan",
		Usage: "Triage SARIF findings using LLM",
		Flags: flags,
		Action: func(ctx context.Context, c *cli.Command) error {
			llmClient, err := llmCfg.Configure(ctx)
			if err != nil {
				return goerr.Wrap(err, "failed to configure LLM")
			}

			repository := fs.New(dataDir)

			uc := usecase.New(
				usecase.WithLLMClient(llmClient),
				usecase.WithRepository(repository),
				usecase.WithTools(
					tool.NewReadFile(repo),
					tool.NewFind(repo),
					tool.NewGrep(repo),
				),
			)

			session, err := uc.Scan(ctx, sarif, repo)
			if err != nil {
				return goerr.Wrap(err, "scan failed")
			}

			out, err := json.MarshalIndent(session, "", "  ")
			if err != nil {
				return goerr.Wrap(err, "failed to marshal output")
			}

			fmt.Println(string(out))
			return nil
		},
	}
}
