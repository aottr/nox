package main

import (
	"log"
	"os"
    "context"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/constants"
	"github.com/aottr/nox/internal/process"
	"github.com/urfave/cli/v3"
)

func main() {

	var configPath string

	cmd := &cli.Command{
		Name:  "nox",
		Usage: "Manage and decrypt app secrets",
		Flags: []cli.Flag{
            &cli.StringFlag{
                Name:        "config",
                Value:       constants.DefaultConfigPath,
                Usage:       "path to config file",
                Destination: &configPath,
            },
        },
        Commands: []*cli.Command{
            {
                Name:    "run",
                Aliases: []string{"r"},
                Usage:   "Fetch, decrypt, and process app secrets",
                Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg, err := config.Load(configPath)
					if err != nil {
						log.Fatalf("failed to load config: %v", err)
					}
                    return process.ProcessApps(cfg)
                },
            },
            {
                Name:    "validate",
                Aliases: []string{"v"},
                Usage:   "Validate configuration and secret integrity",
                Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg, err := config.Load(configPath)
					if err != nil {
						log.Fatalf("failed to load config: %v", err)
					}
                    return process.Validate(cfg)
                },
            },
        },
    }

    if err := cmd.Run(context.Background(), os.Args); err != nil {
        log.Fatal(err)
    }
}
