package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"filippo.io/age"
	"github.com/aottr/nox/internal/crypto"
	"github.com/aottr/nox/internal/state"
)

type RuntimeOptions struct {
	ConfigPath    string
	StatePath     string
	IdentityPaths []string
	DryRun        bool
	Force         bool
	Verbose       bool
	AppName       string
}

type RuntimeContext struct {
	Config     *Config
	State      *state.State
	Identities []age.Identity
	App        *string
	Logger     *log.Logger
	DryRun     bool
	Force      bool
	Verbose    bool
}

func BuildRuntimeContext(opts RuntimeOptions) (*RuntimeContext, error) {

	cfg, err := Load(opts.ConfigPath)
	if err != nil {
		return nil, err
	}

	if opts.StatePath != "" {
		state.SetPath(opts.StatePath)
	}
	st, err := state.Load()
	if err != nil {
		return nil, err
	}

	identityPaths := opts.IdentityPaths
	if len(identityPaths) == 0 {
		// try single identity file first
		if cfg.Age.Identity != "" {
			identityPaths = []string{cfg.Age.Identity}
		} else if len(cfg.Age.Identities) > 0 {
			identityPaths = cfg.Age.Identities
		} else {
			return nil, fmt.Errorf("no age identites found")
		}
	}
	ids, err := crypto.LoadAgeIdentitiesFromPaths(identityPaths)
	if err != nil {
		return nil, err
	}

	var app *string
	if opts.AppName != "" {
		if _, exists := cfg.Apps[opts.AppName]; exists {
			app = &opts.AppName
		} else {
			return nil, fmt.Errorf("app '%s' not found in configuration", opts.AppName)
		}
	}

	var logger *log.Logger
	if opts.Verbose {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		logger = log.New(io.Discard, "", 0)
	}

	return &RuntimeContext{
		Config:     cfg,
		State:      st,
		Identities: ids,
		App:        app,
		Logger:     logger,
		DryRun:     opts.DryRun,
		Force:      opts.Force,
		Verbose:    opts.Verbose,
	}, nil
}
