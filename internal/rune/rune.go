package rune

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Stage string

const (
	StageDraft    Stage = "draft"
	StageReviewed Stage = "reviewed"
	StageStable   Stage = "stable"
)

type Rune struct {
	Name             string   `toml:"name"              json:"name"`
	Description      string   `toml:"description"       json:"description"`
	Version          string   `toml:"version"           json:"version"`
	Stage            Stage    `toml:"stage"             json:"stage"`
	Runtime          string   `toml:"runtime,omitempty" json:"runtime,omitempty"`
	Path             string   `toml:"path,omitempty"    json:"path,omitempty"`
	Inputs           []string `toml:"inputs,omitempty"           json:"inputs,omitempty"`
	Outputs          []string `toml:"outputs,omitempty"          json:"outputs,omitempty"`
	EventsPublished  []string `toml:"events_published,omitempty" json:"events_published,omitempty"`
	EventsSubscribed []string `toml:"events_subscribed,omitempty" json:"events_subscribed,omitempty"`
	Dependencies     []string `toml:"dependencies,omitempty"     json:"dependencies,omitempty"`
	Requirements     []string `toml:"requirements,omitempty"     json:"requirements,omitempty"`
	Config           []string `toml:"config,omitempty"           json:"config,omitempty"`
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
	if r.Stage == "" {
		r.Stage = StageDraft
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

func (s *Store) Promote(name string) (*Rune, error) {
	r, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	switch r.Stage {
	case StageDraft:
		r.Stage = StageReviewed
	case StageReviewed:
		r.Stage = StageStable
	case StageStable:
		return nil, fmt.Errorf("rune %q is already stable", name)
	default:
		return nil, fmt.Errorf("unknown stage %q", r.Stage)
	}
	if err := write(s.filePath(name), *r); err != nil {
		return nil, err
	}
	return r, nil
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
