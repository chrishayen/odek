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
	TokenEnv  string `toml:"token_env,omitempty"` // claude-pro: env var holding CLAUDE_CODE_OAUTH_TOKEN (from `claude setup-token`)
}

type Config struct {
	Agents map[string]Agent `toml:"agents"`
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
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	for name, agent := range cfg.Agents {
		if !validTypes[agent.Type] {
			return nil, fmt.Errorf("agent %q: unknown type %q (valid: claude-api, claude-pro, docker)", name, agent.Type)
		}
	}
	return &cfg, nil
}
