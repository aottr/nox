package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/constants"
	"github.com/aottr/nox/internal/processor"
	"github.com/urfave/cli/v3"
)

func main() {

	var configPath string
	var statePath string
	var identityPath string
	var appName string
	var dryRun bool
	var force bool
	var verbose bool
	var inputPath string
	var outputPath string

	cmd := &cli.Command{
		Name:                  "nox",
		Usage:                 "Manage and decrypt app secrets",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       constants.DefaultConfigPath,
				Usage:       "path to config file",
				Destination: &configPath,
			},
			&cli.StringFlag{
				Name:        "state",
				Value:       constants.DefaultStatePath,
				Usage:       "path to state file",
				Destination: &statePath,
			},
			&cli.StringFlag{
				Name:        "identity",
				Usage:       "path to age identity file",
				Destination: &identityPath,
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Value:       false,
				Usage:       "print verbose output",
				Destination: &verbose,
			},
		},
		Commands: []*cli.Command{
			// {
			// 	Name:    "run",
			// 	Aliases: []string{"r"},
			// 	Usage:   "Fetch, decrypt, and process app secrets",
			// 	Action: func(ctx context.Context, cmd *cli.Command) error {
			// 		cfg, err := config.Load(configPath)
			// 		if err != nil {
			// 			log.Fatalf("failed to load config: %v", err)
			// 		}
			// 		return processor.ProcessApps(cfg)
			// 	},
			// },
			{
				Name:    "export",
				Aliases: []string{"e"},
				Usage:   "Export all secrets to a single file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "app",
						Aliases:     []string{"a"},
						Usage:       "app to export secrets for",
						Destination: &appName,
					},
					&cli.BoolFlag{
						Name:        "dry-run",
						Aliases:     []string{"d"},
						Value:       false,
						Usage:       "only print what would be exported",
						Destination: &dryRun,
					},
					&cli.BoolFlag{
						Name:        "force",
						Aliases:     []string{"f"},
						Value:       false,
						Usage:       "ignore state file",
						Destination: &force,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rtx, err := config.BuildRuntimeContext(config.RuntimeOptions{
						ConfigPath:   configPath,
						StatePath:    statePath,
						IdentityPath: identityPath,
						DryRun:       dryRun,
						Force:        force,
						AppName:      appName,
						Verbose:      verbose,
					})
					if err != nil {
						log.Fatalf("failed to build runtime context: %v", err)
					}
					if appName != "" {
						return processor.SyncApp(rtx)
					}
					return processor.SyncApps(rtx)
				},
			},
			{
				Name:    "encrypt",
				Aliases: []string{"enc"},
				Usage:   "Encrypt a file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Usage:       "path to input file",
						Destination: &inputPath,
					},
					&cli.StringFlag{
						Name:        "output",
						Usage:       "path to output file",
						Destination: &outputPath,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("encrypting file", inputPath)
					fmt.Println("writing to", outputPath)
					return nil
				},
			},
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "Validate configuration and secret integrity",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rtx, err := config.BuildRuntimeContext(config.RuntimeOptions{
						ConfigPath:   configPath,
						StatePath:    statePath,
						IdentityPath: identityPath,
					})
					if err != nil {
						log.Fatalf("failed to build runtime context: %v", err)
					}
					return processor.ValidateConfig(rtx.Config)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
