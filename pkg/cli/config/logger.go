package config

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/m-mizutani/goerr/v2"
	"github.com/secmon-lab/queen/pkg/utils/logging"
	"github.com/urfave/cli/v3"
)

type Logger struct {
	level      string
	format     string
	output     string
	stacktrace bool
}

func (x *Logger) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Category:    "Logging",
			Aliases:     []string{"l"},
			Sources:     cli.EnvVars("QUEEN_LOG_LEVEL"),
			Usage:       "Set log level [debug|info|warn|error]",
			Value:       "info",
			Destination: &x.level,
		},
		&cli.StringFlag{
			Name:        "log-format",
			Category:    "Logging",
			Sources:     cli.EnvVars("QUEEN_LOG_FORMAT"),
			Usage:       "Set log format [console|json]",
			Value:       "",
			Destination: &x.format,
		},
		&cli.StringFlag{
			Name:        "log-output",
			Category:    "Logging",
			Sources:     cli.EnvVars("QUEEN_LOG_OUTPUT"),
			Usage:       "Set log output (stderr, stdout, or file path)",
			Value:       "stderr",
			Destination: &x.output,
		},
		&cli.BoolFlag{
			Name:        "log-stacktrace",
			Category:    "Logging",
			Sources:     cli.EnvVars("QUEEN_LOG_STACKTRACE"),
			Usage:       "Show stacktrace in error logs",
			Destination: &x.stacktrace,
			Value:       true,
		},
	}
}

func (x *Logger) Configure() (*slog.Logger, func(), error) {
	closer := func() {}

	levelMap := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	level, ok := levelMap[x.level]
	if !ok {
		return nil, closer, goerr.New("invalid log level", goerr.V("level", x.level))
	}

	formatMap := map[string]logging.Format{
		"console": logging.FormatConsole,
		"json":    logging.FormatJSON,
	}

	var format logging.Format
	if x.format == "" {
		term := os.Getenv("TERM")
		if strings.Contains(term, "color") || strings.Contains(term, "xterm") {
			format = logging.FormatConsole
		} else {
			format = logging.FormatJSON
		}
	} else {
		f, ok := formatMap[x.format]
		if !ok {
			return nil, closer, goerr.New("invalid log format", goerr.V("format", x.format))
		}
		format = f
	}

	var output io.Writer
	switch x.output {
	case "stdout", "-":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		f, err := os.OpenFile(filepath.Clean(x.output), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			return nil, closer, goerr.Wrap(err, "failed to open log file", goerr.V("path", x.output))
		}
		output = f
		closer = func() { f.Close() }
	}

	logger := logging.New(output, level, format, x.stacktrace)
	logging.SetDefault(logger)

	return logger, closer, nil
}
