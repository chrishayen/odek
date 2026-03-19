package rune

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Rune struct {
	Name        string `toml:"name"        json:"name"`
	Description string `toml:"description" json:"description"` // English description — this is the spec
	Version     string `toml:"version"     json:"version"`
	Hydrated    bool   `toml:"hydrated"    json:"hydrated"`  // true once code has been generated
	Coverage    float64 `toml:"coverage"   json:"coverage"`  // test coverage %, set after hydration
}

type Store struct {
	registryPath string
}

func NewStore(registryPath string) *Store {
	return &Store{registryPath: registryPath}
}

func (s *Store) dir() string {
	return filepath.Join(s.registryPath, "runes")
}

func (s *Store) filePath(name string) string {
	return filepath.Join(s.dir(), name+".toml")
}

// CodeDir returns the directory where generated code for a rune is stored.
func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.dir(), name)
}

func (s *Store) Create(r Rune) error {
	if err := validate(r); err != nil {
		return err
	}
	if err := os.MkdirAll(s.dir(), 0755); err != nil {
		return fmt.Errorf("creating registry dir: %w", err)
	}
	if _, err := os.Stat(s.filePath(r.Name)); err == nil {
		return fmt.Errorf("rune %q already exists", r.Name)
	}
	if r.Version == "" {
		r.Version = "0.1.0"
	}
	return write(s.filePath(r.Name), r)
}

func (s *Store) Get(name string) (*Rune, error) {
	var r Rune
	if _, err := toml.DecodeFile(s.filePath(name), &r); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("rune %q not found", name)
		}
		return nil, err
	}
	return &r, nil
}

func (s *Store) List() ([]Rune, error) {
	entries, err := os.ReadDir(s.dir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var runes []Rune
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".toml") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".toml")
		r, err := s.Get(name)
		if err != nil {
			continue
		}
		runes = append(runes, *r)
	}
	return runes, nil
}

func (s *Store) Update(r Rune) error {
	if err := validate(r); err != nil {
		return err
	}
	if _, err := os.Stat(s.filePath(r.Name)); os.IsNotExist(err) {
		return fmt.Errorf("rune %q not found", r.Name)
	}
	return write(s.filePath(r.Name), r)
}

func (s *Store) Delete(name string) error {
	p := s.filePath(name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("rune %q not found", name)
	}
	return os.Remove(p)
}

func validate(r Rune) error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if strings.ContainsAny(r.Name, " /\\") {
		return fmt.Errorf("name must be a slug (no spaces or slashes)")
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}

func write(path string, r Rune) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(r)
}
