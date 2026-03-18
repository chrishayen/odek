package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Agent struct {
	Type      string `toml:"type"`
	Model     string `toml:"model,omitempty"`
	APIKeyEnv string `toml:"api_key_env,omitempty"`
	Image     string `toml:"image,omitempty"`
}

type Config struct {
	Agents map[string]Agent `toml:"agents"`
}

func Load() (*Config, error) {
	path := os.Getenv("VALKYRIE_CONFIG")
	if path == "" {
		return nil, fmt.Errorf("VALKYRIE_CONFIG is not set")
	}
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	return &cfg, nil
}
