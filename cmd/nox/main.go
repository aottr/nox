package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aottr/nox/internal/config"
	"github.com/aottr/nox/internal/constants"
	"github.com/aottr/nox/internal/crypto"
	"github.com/aottr/nox/internal/logging"
	"github.com/aottr/nox/internal/processor"
	"github.com/urfave/cli/v3"
)

func main() {

	logging.Init()
	logging.SetLevel("info")
	log := logging.Get()

	var configPath string
	var statePath string
	var identityPaths []string

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
			&cli.StringSliceFlag{
				Name:        "identity",
				Usage:       "path to age identity file",
				Destination: &identityPaths,
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
				Name:    "encrypt",
				Aliases: []string{"enc"},
				Usage:   "Encrypt a file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "input",
						Usage:       "path to input file",
						Aliases:     []string{"i"},
						Value:       constants.StandardInput,
						Destination: &inputPath,
					},
					&cli.StringFlag{
						Name:        "output",
						Usage:       "path to output file",
						Aliases:     []string{"o"},
						Value:       constants.StandardOutput,
						Destination: &outputPath,
					},
					&cli.StringSliceFlag{
						Name:    "recipient",
						Usage:   "age public key of recipient (repeatable)",
						Aliases: []string{"r"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("encrypting file", inputPath)
					fmt.Println("writing to", outputPath)
					recipients, err := crypto.StringsToRecipients(cmd.StringSlice("recipient"))
					if err != nil {
						return err
					}
					out, err := crypto.EncryptFile(inputPath, recipients)
					if err != nil {
						return err
					}
					fmt.Println(string(out))
					// priv, pub, err := crypto.GenerateAndWriteX25519Identity("test.key")
					// if err != nil {
					// 	return err
					// }
					// fmt.Println(priv)
					// fmt.Println(pub)

					return nil
				},
			},
			{
				Name: "generate",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "output",
						Usage:       "path to output file",
						Aliases:     []string{"o"},
						Value:       constants.StandardOutput,
						Destination: &outputPath,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					priv, pub, err := crypto.GenerateIdentity(cmd.String("output"))
					if err != nil {
						return err
					}
					fmt.Println(priv)
					fmt.Println(pub)
					return nil
				},
			},
			{
				Name:    "decrypt",
				Aliases: []string{"d"},
				Usage:   "Decrypts all secrets of one or all apps",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "app",
						Aliases: []string{"a"},
						Usage:   "app to decrypt secrets for",
					},
					&cli.BoolFlag{
						Name:        "dry-run",
						Aliases:     []string{"d"},
						Value:       false,
						Usage:       "only print what would be decrypted",
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
						ConfigPath:    configPath,
						StatePath:     statePath,
						IdentityPaths: identityPaths,
						DryRun:        dryRun,
						Force:         force,
						AppName:       cmd.String("app"),
						Verbose:       verbose,
					})
					if err != nil {
						log.Error("failed to build runtime context", "error", err.Error())
					}
					if cmd.String("app") != "" {
						return processor.SyncApp(rtx)
					}
					return processor.SyncApps(rtx)
				},
			},
			{
				Name:    "validate",
				Aliases: []string{"v"},
				Usage:   "Validate configuration and secret integrity",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					rtx, err := config.BuildRuntimeContext(config.RuntimeOptions{
						ConfigPath:    configPath,
						StatePath:     statePath,
						IdentityPaths: identityPaths,
					})
					if err != nil {
						log.Error("failed to build runtime context", "error", err.Error())
					}
					return processor.ValidateConfig(rtx.Config)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Error("failed to run command", "error", err.Error())
	}
}
