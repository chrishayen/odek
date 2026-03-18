package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DefaultSandbox string `toml:"default_sandbox"`
	RegistryPath   string `toml:"registry_path"`
}

func Load(path string) (*Config, error) {
	cfg := &Config{
		RegistryPath: "./registry",
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
