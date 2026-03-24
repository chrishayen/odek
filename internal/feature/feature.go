package feature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
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
}

func NewStore(registryPath string) *Store {
	return &Store{registryPath: registryPath}
}

func (s *Store) filePath(name string) string {
	return filepath.Join(s.registryPath, name, "feature.md")
}

func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.registryPath, name, "_composed")
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
	f := parseFrontmatter(string(data))
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
	body := stripFrontmatter(string(data))

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

// parseFrontmatter extracts only YAML frontmatter fields into a Feature.
func parseFrontmatter(content string) Feature {
	var f Feature
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---")
		if end != -1 {
			fm := content[4 : 4+end]
			_ = yaml.Unmarshal([]byte(fm), &f)
		}
	}
	return f
}

// stripFrontmatter returns everything after the closing --- delimiter.
func stripFrontmatter(content string) string {
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---")
		if end != -1 {
			return content[4+end+4:]
		}
	}
	return content
}
