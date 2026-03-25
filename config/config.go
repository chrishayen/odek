package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Agent struct {
	Type     string `toml:"type"`
	Model    string `toml:"model,omitempty"`
	Image    string `toml:"image,omitempty"`
	Token    string `toml:"token,omitempty"`
	TokenEnv string `toml:"token_env,omitempty"`
	Sandbox  bool   `toml:"sandbox,omitempty"`
}

type Config struct {
	Project      string `toml:"project"`
	Language     string `toml:"language"`
	RegistryPath string `toml:"registry_path"`
	OutputPath   string `toml:"output_path"`
	Agent        Agent  `toml:"agent"`
}

var supportedLanguages = map[string]bool{
	"go": true,
}

var validTypes = map[string]bool{
	"claude-sub": true,
	"mock":       true, // for testing only
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
		RegistryPath: filepath.Join(root, "runes"),
		OutputPath:   filepath.Join(root, "src"),
		Agent: Agent{
			Type: "claude-sub",
		},
	}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	if cfg.Project == "" {
		return nil, fmt.Errorf("project name is required in %s", path)
	}

	if cfg.Language == "" {
		cfg.Language = "go"
	}
	if !supportedLanguages[cfg.Language] {
		return nil, fmt.Errorf("unsupported language %q (supported: go)", cfg.Language)
	}

	if !validTypes[cfg.Agent.Type] {
		return nil, fmt.Errorf("agent: unknown type %q (valid: claude-sub)", cfg.Agent.Type)
	}

	return cfg, nil
}

// ClaudeToken reads the OAuth access token from ~/.claude/.credentials.json.
// Returns empty string if the file doesn't exist or can't be parsed.
func ClaudeToken() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(home, ".claude", ".credentials.json"))
	if err != nil {
		return ""
	}
	var creds struct {
		ClaudeAiOauth struct {
			AccessToken string `json:"accessToken"`
		} `json:"claudeAiOauth"`
	}
	if json.Unmarshal(data, &creds) != nil {
		return ""
	}
	return creds.ClaudeAiOauth.AccessToken
}
