package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Agent struct {
	Model    string `toml:"model,omitempty"`
	Token    string `toml:"token,omitempty"`
	TokenEnv string `toml:"token_env,omitempty"`
	Mock     bool   `toml:"mock,omitempty"`
}

type Server struct {
	Port int `toml:"port,omitempty"`
}

type Config struct {
	Project      string `toml:"project"`
	Language     string `toml:"language"`
	RegistryPath string `toml:"registry_path"`
	OutputPath   string `toml:"output_path"`
	Concurrency  int    `toml:"concurrency,omitempty"`
	Agent        Agent  `toml:"agent"`
	Server       Server `toml:"server"`
}

var supportedLanguages = map[string]bool{
	"go": true,
	"ts": true,
	"py": true,
}

// ResolveToken returns the OAuth token using this precedence:
// 1. token field in config (literal value)
// 2. token_env field in config (named env var)
// 3. CLAUDE_CODE_OAUTH_TOKEN env var
// 4. Parse from ~/.claude/.credentials.json
func (a Agent) ResolveToken() string {
	if a.Token != "" {
		return a.Token
	}
	if a.TokenEnv != "" {
		if v := os.Getenv(a.TokenEnv); v != "" {
			return v
		}
	}
	if v := os.Getenv("CLAUDE_CODE_OAUTH_TOKEN"); v != "" {
		return v
	}
	return ClaudeToken()
}

// FindRoot walks up from cwd (or ODEK_PROJECT_DIR) looking for odek.toml.
func FindRoot() (string, error) {
	start := os.Getenv("ODEK_PROJECT_DIR")
	if start == "" {
		var err error
		start, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	dir := start
	for {
		path := filepath.Join(dir, "odek.toml")
		if _, err := os.Stat(path); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("odek.toml not found — run 'odek init' to create a project")
		}
		dir = parent
	}
}

func Load() (*Config, error) {
	root, err := FindRoot()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, "odek.toml")
	cfg := &Config{
		RegistryPath: root,
		OutputPath:   filepath.Join(root, "src"),
		Concurrency:  50,
		Server: Server{
			Port: 8319,
		},
	}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}

	if cfg.Language == "" {
		cfg.Language = "go"
	}
	if !supportedLanguages[cfg.Language] {
		return nil, fmt.Errorf("unsupported language %q (supported: go, ts, py)", cfg.Language)
	}

	return cfg, nil
}

// ClaudeToken reads the OAuth access token from ~/.claude/.credentials.json.
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

// LangExtension returns the file extension for a language identifier.
func LangExtension(lang string) string {
	switch lang {
	case "go":
		return ".go"
	case "py":
		return ".py"
	default:
		return ".ts"
	}
}

// LangName returns the human-readable name for a language identifier.
func LangName(lang string) string {
	switch lang {
	case "go":
		return "Go"
	case "py":
		return "Python"
	default:
		return "TypeScript (Node.js with node:* built-ins)"
	}
}
