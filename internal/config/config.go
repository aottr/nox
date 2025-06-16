package config

import (
	"os"

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

type AppConfig struct {
	Repo   string   `yaml:"repo,omitempty"`
	Branch string   `yaml:"branch"`
	Files  []FileConfig `yaml:"files"`
}

type Config struct {
	Interval    string               `yaml:"interval"`
	AgeKeyPath  string               `yaml:"ageKeyPath"`
	StatePath   string               `yaml:"statePath"`
	DefaultRepo string               `yaml:"defaultRepo"`
	Secrets     []SecretMapping      `yaml:"secrets"`
	Apps        map[string]AppConfig `yaml:"apps"`
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

	return &cfg, nil
}
