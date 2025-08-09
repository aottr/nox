package main

import (
	"github.com/aottr/nox/internal/constants"
)

type InitOptions struct {
	ConfigPath       string
	IdentityPath     string
	GenerateIdentity *bool
}

func RunInit(opts *InitOptions) error {
	if opts.ConfigPath == "" {
		opts.ConfigPath = constants.DefaultConfigPath
	}

	return nil
}
