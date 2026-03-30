package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type App struct {
	Name        string   `json:"name"         yaml:"-"`
	Description string   `json:"description"  yaml:"-"`
	Version     string   `json:"version"      yaml:"version"`
	Status      string   `json:"status"       yaml:"status"`
	Features    []string `json:"features"     yaml:"features"`
	EntryPoint  string   `json:"entry_point"  yaml:"entry_point"`
	Hydrated    bool     `json:"hydrated"     yaml:"hydrated"`
	Coverage    float64  `json:"coverage"     yaml:"coverage"`
	Raw         string   `json:"raw"          yaml:"-"`
}

type Store struct {
	registryPath string
	outputPath   string
}

func NewStore(registryPath, outputPath string) *Store {
	return &Store{registryPath: registryPath, outputPath: outputPath}
}

func (s *Store) filePath(name string) string {
	return filepath.Join(s.registryPath, name, "app.md")
}

func (s *Store) OutputPath() string {
	return s.outputPath
}

func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.outputPath, name)
}

func (s *Store) ReadRaw(name string) (string, error) {
	data, err := os.ReadFile(s.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("app %q not found", name)
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
		return fmt.Errorf("app name must be a single slug (no spaces, slashes, or backslashes)")
	}
	if body == "" {
		return fmt.Errorf("body is required")
	}

	p := s.filePath(name)
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("app %q already exists", name)
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("creating app dir: %w", err)
	}

	content := fmt.Sprintf("---\nversion: 0.1.0\nstatus: draft\nhydrated: false\ncoverage: -1\nfeatures: []\nentry_point: \"\"\n---\n\n# %s\n\n%s\n", name, body)
	return os.WriteFile(p, []byte(content), 0644)
}

func (s *Store) Get(name string) (*App, error) {
	data, err := os.ReadFile(s.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("app %q not found", name)
		}
		return nil, err
	}
	a := parseFrontmatter(string(data))
	a.Name = name
	a.Raw = string(data)
	return &a, nil
}

func (s *Store) List() ([]App, error) {
	base := s.registryPath
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, nil
	}
	var apps []App
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
		if filepath.Join(base, e.Name()) == s.outputPath {
			continue
		}
		fp := filepath.Join(base, e.Name(), "app.md")
		if _, err := os.Stat(fp); os.IsNotExist(err) {
			continue
		}
		a, err := s.Get(e.Name())
		if err != nil {
			continue
		}
		apps = append(apps, *a)
	}
	return apps, nil
}

func (s *Store) Update(a App) error {
	if a.Name == "" {
		return fmt.Errorf("name is required")
	}
	p := s.filePath(a.Name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("app %q not found", a.Name)
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	body := stripFrontmatter(string(data))

	var sb strings.Builder
	fmt.Fprintln(&sb, "---")
	fmt.Fprintf(&sb, "version: %s\n", a.Version)
	fmt.Fprintf(&sb, "status: %s\n", a.Status)
	fmt.Fprintf(&sb, "hydrated: %t\n", a.Hydrated)
	fmt.Fprintf(&sb, "coverage: %g\n", a.Coverage)
	if len(a.Features) > 0 {
		fmt.Fprintln(&sb, "features:")
		for _, f := range a.Features {
			fmt.Fprintf(&sb, "  - %s\n", f)
		}
	} else {
		fmt.Fprintln(&sb, "features: []")
	}
	if a.EntryPoint != "" {
		fmt.Fprintf(&sb, "entry_point: %s\n", a.EntryPoint)
	} else {
		fmt.Fprintln(&sb, "entry_point: \"\"")
	}
	fmt.Fprintln(&sb, "---")
	sb.WriteString(body)

	return os.WriteFile(p, []byte(sb.String()), 0644)
}

func (s *Store) Delete(name string) error {
	p := s.filePath(name)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("app %q not found", name)
	}
	return os.Remove(p)
}

func parseFrontmatter(content string) App {
	var a App
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---")
		if end != -1 {
			fm := content[4 : 4+end]
			_ = yaml.Unmarshal([]byte(fm), &a)
		}
	}
	return a
}

func stripFrontmatter(content string) string {
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---")
		if end != -1 {
			return content[4+end+4:]
		}
	}
	return content
}
