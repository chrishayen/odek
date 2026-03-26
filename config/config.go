package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Server struct {
	URL      string `toml:"url"`
	Token    string `toml:"token,omitempty"`
	TokenEnv string `toml:"token_env,omitempty"`
}

type Config struct {
	Project string `toml:"project"`
	Server  Server `toml:"server"`
}

// FindRoot walks up from cwd (or VALKYRIE_PROJECT_DIR) looking for valkyrie.toml.
func FindRoot() (string, error) {
	start := os.Getenv("VALKYRIE_PROJECT_DIR")
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	dir := start
	for {
		path := filepath.Join(dir, "valkyrie.toml")
		if _, err := os.Stat(path); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("valkyrie.toml not found — run 'valkyrie init' to create a project")
		}
		dir = parent
	}
}

func Load() (*Config, error) {
	root, err := FindRoot()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, "valkyrie.toml")
	cfg := &Config{
		Server: Server{
			URL:      "http://localhost:7777",
			TokenEnv: "VALKYRIE_TOKEN",
		},
	}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	if cfg.Project == "" {
		return nil, fmt.Errorf("project name is required in %s", path)
	}

	if cfg.Server.URL == "" {
		return nil, fmt.Errorf("server.url is required in %s", path)
	}

	return cfg, nil
}

// ResolveToken returns the server token from config or environment.
func (c *Config) ResolveToken() string {
	if c.Server.Token != "" {
		return c.Server.Token
	}
	if c.Server.TokenEnv != "" {
		return os.Getenv(c.Server.TokenEnv)
	}
	return ""
}
