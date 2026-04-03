package feature

import (
	"fmt"
	"path/filepath"

	runepkg "github.com/chrishayen/odek/internal/rune"
)

// Feature is a view of a top-level rune package.
type Feature struct {
	Name     string  `json:"name"`
	Version  string  `json:"version"`
	Hydrated bool    `json:"hydrated"`
	Coverage float64 `json:"coverage"`
}

type Store struct {
	runeStore  *runepkg.Store
	outputPath string
}

func NewStore(runeStore *runepkg.Store, outputPath string) *Store {
	return &Store{runeStore: runeStore, outputPath: outputPath}
}

func (s *Store) OutputPath() string {
	return s.outputPath
}

func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.outputPath, name)
}

func (s *Store) List() ([]Feature, error) {
	pkgs, err := s.runeStore.TopLevelPackages()
	if err != nil {
		return nil, err
	}
	features := make([]Feature, len(pkgs))
	for i, r := range pkgs {
		features[i] = toFeature(r)
	}
	return features, nil
}

func (s *Store) Get(name string) (*Feature, error) {
	r, err := s.runeStore.Get(name)
	if err != nil {
		return nil, fmt.Errorf("feature %q not found", name)
	}
	f := toFeature(*r)
	return &f, nil
}

func toFeature(r runepkg.Rune) Feature {
	return Feature{
		Name:     r.Name,
		Version:  r.Version.String(),
		Hydrated: r.Hydrated,
		Coverage: r.Coverage,
	}
}
