package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Sandbox struct {
	Name      string `toml:"name"`
	Type      string `toml:"type"`
	Model     string `toml:"model,omitempty"`
	APIKeyEnv string `toml:"api_key_env,omitempty"`
	Image     string `toml:"image,omitempty"`
}

func dir(registryPath string) string {
	return filepath.Join(registryPath, "sandboxes")
}

func path(registryPath, name string) string {
	return filepath.Join(dir(registryPath), name+".toml")
}

func Create(registryPath string, s Sandbox) error {
	if err := os.MkdirAll(dir(registryPath), 0755); err != nil {
		return fmt.Errorf("creating sandbox dir: %w", err)
	}
	p := path(registryPath, s.Name)
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("sandbox %q already exists", s.Name)
	}
	f, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(s)
}

func Get(registryPath, name string) (*Sandbox, error) {
	var s Sandbox
	if _, err := toml.DecodeFile(path(registryPath, name), &s); err != nil {
		return nil, fmt.Errorf("sandbox %q not found", name)
	}
	return &s, nil
}

func List(registryPath string) ([]Sandbox, error) {
	entries, err := os.ReadDir(dir(registryPath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var sandboxes []Sandbox
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".toml")
		s, err := Get(registryPath, name)
		if err != nil {
			continue
		}
		sandboxes = append(sandboxes, *s)
	}
	return sandboxes, nil
}

func Delete(registryPath, name string) error {
	p := path(registryPath, name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("sandbox %q not found", name)
	}
	return os.Remove(p)
}
