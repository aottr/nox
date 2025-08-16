package config

import (
	"fmt"

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
	App        string
	DryRun     bool
	Force      bool
}

func BuildRuntimeCtxFromConfig(config *Config) (*RuntimeContext, error) {

	if config.StatePath != "" {
		state.SetPath(config.StatePath)
	}
	st, err := state.Load()
	if err != nil {
		return nil, err
	}

	var identityPaths []string
	// try single identity file first
	if config.Age.Identity != "" {
		identityPaths = []string{config.Age.Identity}
	} else if len(config.Age.Identities) > 0 {
		identityPaths = config.Age.Identities
	} else {
		return nil, fmt.Errorf("no age identites found")
	}
	ids, err := crypto.LoadAgeIdentitiesFromPaths(identityPaths)
	if err != nil {
		return nil, err
	}

	return &RuntimeContext{
		Config:     config,
		State:      st,
		Identities: ids,
		DryRun:     false,
		Force:      false,
	}, nil
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

	var app string
	if opts.AppName != "" {
		if _, exists := cfg.Apps[opts.AppName]; exists {
			app = opts.AppName
		} else {
			return nil, fmt.Errorf("app '%s' not found in configuration", opts.AppName)
		}
	}

	return &RuntimeContext{
		Config:     cfg,
		State:      st,
		Identities: ids,
		App:        app,
		DryRun:     opts.DryRun,
		Force:      opts.Force,
	}, nil
}
