package feature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chrishayen/odek/internal/frontmatter"
)

type Feature struct {
	Name     string  `json:"name"     yaml:"-"`
	Version  string  `json:"version"  yaml:"version"`
	Status   string  `json:"status"   yaml:"status"`
	Hydrated bool    `json:"hydrated" yaml:"hydrated"`
	Coverage float64 `json:"coverage" yaml:"coverage"`
	Raw      string  `json:"raw"      yaml:"-"`
}

type Store struct {
	registryPath string
	outputPath   string
}

func NewStore(registryPath, outputPath string) *Store {
	return &Store{registryPath: registryPath, outputPath: outputPath}
}

func (s *Store) filePath(name string) string {
	return filepath.Join(s.registryPath, name, "feature.md")
}

// OutputPath returns the root output directory for generated code.
func (s *Store) OutputPath() string {
	return s.outputPath
}

func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.outputPath, name)
}

// ReadRaw returns the raw file content for passing to agents.
func (s *Store) ReadRaw(name string) (string, error) {
	data, err := os.ReadFile(s.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("feature %q not found", name)
		}
		return "", err
	}
	return string(data), nil
}

func (s *Store) Create(name, body string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if strings.ContainsAny(name, " \\/") {
		return fmt.Errorf("feature name must be a single slug (no spaces, slashes, or backslashes)")
	}
	if body == "" {
		return fmt.Errorf("body is required")
	}

	p := s.filePath(name)
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("feature %q already exists", name)
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("creating feature dir: %w", err)
	}

	content := fmt.Sprintf("---\nversion: 0.1.0\nstatus: draft\nhydrated: false\ncoverage: -1\n---\n\n# %s\n\n%s\n", name, body)
	return os.WriteFile(p, []byte(content), 0644)
}

func (s *Store) Get(name string) (*Feature, error) {
	data, err := os.ReadFile(s.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("feature %q not found", name)
		}
		return nil, err
	}
	var f Feature
	frontmatter.Parse(string(data), &f)
	f.Name = name
	f.Raw = string(data)
	return &f, nil
}

func (s *Store) List() ([]Feature, error) {
	base := s.registryPath
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, nil
	}
	var features []Feature
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if e.Name() == "draft" {
			continue
		}
		if filepath.Join(base, e.Name()) == s.outputPath {
			continue
		}
		fp := filepath.Join(base, e.Name(), "feature.md")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			continue
		}
		f, err := s.Get(e.Name())
		if err != nil {
			continue
		}
		features = append(features, *f)
	}
	return features, nil
}

// Update replaces the frontmatter in the feature file, preserving the body.
func (s *Store) Update(f Feature) error {
	if f.Name == "" {
		return fmt.Errorf("name is required")
	}
	p := s.filePath(f.Name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("feature %q not found", f.Name)
	}

	// Read existing file to preserve body
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	body := frontmatter.Strip(string(data))

	// Write new frontmatter + preserved body
	var sb strings.Builder
	fmt.Fprintln(&sb, "---")
	fmt.Fprintf(&sb, "version: %s\n", f.Version)
	fmt.Fprintf(&sb, "status: %s\n", f.Status)
	fmt.Fprintf(&sb, "hydrated: %t\n", f.Hydrated)
	fmt.Fprintf(&sb, "coverage: %g\n", f.Coverage)
	fmt.Fprintln(&sb, "---")
	sb.WriteString(body)

	return os.WriteFile(p, []byte(sb.String()), 0644)
}

func (s *Store) Delete(name string) error {
	p := s.filePath(name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("feature %q not found", name)
	}
	return os.Remove(p)
}

