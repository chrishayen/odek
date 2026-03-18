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
	Token     string `toml:"token,omitempty"`
	TokenEnv  string `toml:"token_env,omitempty"`
}

type Auth struct {
	Disabled bool `toml:"disabled"` // set true to acknowledge no auth (local use)
}

type Config struct {
	RegistryPath string           `toml:"registry_path"`
	Auth         Auth             `toml:"auth"`
	Agents       map[string]Agent `toml:"agents"`
}

var validTypes = map[string]bool{
	"claude-api": true,
	"claude-pro": true,
	"docker":     true,
}

func Load() (*Config, error) {
	path := os.Getenv("VALKYRIE_CONFIG")
	if path == "" {
		return nil, fmt.Errorf("VALKYRIE_CONFIG is not set")
	}
	cfg := &Config{
		RegistryPath: "./registry",
	}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	for name, agent := range cfg.Agents {
		if !validTypes[agent.Type] {
			return nil, fmt.Errorf("agent %q: unknown type %q (valid: claude-api, claude-pro, docker)", name, agent.Type)
		}
	}
	return cfg, nil
}
