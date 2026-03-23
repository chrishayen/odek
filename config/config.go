package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Dir returns the Valkyrie config directory (~/.config/valkyrie).
func Dir() string {
	if d := os.Getenv("VALKYRIE_CONFIG_DIR"); d != "" {
		return d
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "valkyrie")
}

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
	"claude-max": true,
	"docker":     true,
	"mock":       true, // for testing only
}

func Load() (*Config, error) {
	dir := Dir()
	path := filepath.Join(dir, "config.toml")

	cfg := &Config{
		RegistryPath: filepath.Join(dir, "registry"),
	}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config not found at %s — run 'valkyrie init' or create it manually", path)
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	for name, agent := range cfg.Agents {
		if !validTypes[agent.Type] {
			return nil, fmt.Errorf("agent %q: unknown type %q (valid: claude-api, claude-max, docker)", name, agent.Type)
		}
	}
	return cfg, nil
}
