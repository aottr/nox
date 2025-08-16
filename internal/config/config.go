package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type SecretMapping struct {
	EncryptedPath string `yaml:"encrypted"`
	OutputPath    string `yaml:"output"`
}

type FileConfig struct {
	Path   string `yaml:"path"`
	Output string `yaml:"output,omitempty"`
}

type GitConfig struct {
	Repo   string `yaml:"repo"`
	Branch string `yaml:"branch"`
}

func (g GitConfig) IsValid() bool {
	return g.Repo != "" && g.Branch != ""
}

type AppConfig struct {
	GitConfig GitConfig    `yaml:"git,omitempty"`
	Files     []FileConfig `yaml:"files"`
}

type AgeConfig struct {
	Identity   string   `yaml:"identity"`
	Identities []string `yaml:"identities,omitempty"`
	Recipients []string `yaml:"recipients,omitempty"`
}

type Config struct {
	Interval       time.Duration        `yaml:"-"`
	IntervalString string               `yaml:"interval"`
	Age            AgeConfig            `yaml:"age"`
	StatePath      string               `yaml:"statePath"`
	GitConfig      GitConfig            `yaml:"git"`
	Apps           map[string]AppConfig `yaml:"apps"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// validate config
	// must have at least one git config
	hasAnyGit := cfg.GitConfig.IsValid()
	if !hasAnyGit {
		for _, app := range cfg.Apps {
			if app.GitConfig.IsValid() {
				hasAnyGit = true
				break
			}
		}
	}
	if !hasAnyGit {
		return nil, fmt.Errorf("no git configuration found: set either top-level git or app-specific git")
	}

	// validate interval
	cfg.Interval, err = time.ParseDuration(cfg.IntervalString)
	if err != nil {
		return nil, fmt.Errorf("invalid interval: %w", err)
	}
	return &cfg, nil
}

func InitConfig(path string) error {

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("config file already exists")
	}

	cfg := Config{
		IntervalString: "10m",
		StatePath:      ".nox-state.json",
		GitConfig: GitConfig{
			Repo:   "",
			Branch: "main",
		},
	}
	cfgYaml, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	return os.WriteFile(path, cfgYaml, 0600)
}
