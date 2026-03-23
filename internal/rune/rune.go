package rune

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Rune struct {
	Name          string   `json:"name"                    yaml:"-"`
	Description   string   `json:"description"             yaml:"-"`
	Behavior      string   `json:"behavior,omitempty"       yaml:"-"`
	PositiveTests []string `json:"positive_tests,omitempty" yaml:"-"`
	NegativeTests []string `json:"negative_tests,omitempty" yaml:"-"`
	Version       string   `json:"version"                 yaml:"version"`
	Hydrated      bool     `json:"hydrated"                yaml:"hydrated"`
	Coverage      float64  `json:"coverage"                yaml:"coverage"`
}

type Store struct {
	registryPath string
}

func NewStore(registryPath string) *Store {
	return &Store{registryPath: registryPath}
}

func (s *Store) dir() string {
	return s.registryPath
}

func (s *Store) filePath(name string) string {
	return filepath.Join(s.dir(), name+".md")
}

// CodeDir returns the directory where generated code for a rune is stored.
func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.dir(), name)
}

func (s *Store) Create(r Rune) error {
	if err := validate(r); err != nil {
		return err
	}
	p := s.filePath(r.Name)
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("creating runes dir: %w", err)
	}
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("rune %q already exists", r.Name)
	}
	if r.Version == "" {
		r.Version = "0.1.0"
	}
	return write(p, r)
}

func (s *Store) Get(name string) (*Rune, error) {
	data, err := os.ReadFile(s.filePath(name))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("rune %q not found", name)
		}
		return nil, err
	}
	r, err := parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing rune %q: %w", name, err)
	}
	return r, nil
}

func (s *Store) List() ([]Rune, error) {
	base := s.dir()
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, nil
	}
	var runes []Rune
	err := filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		rel, _ := filepath.Rel(base, path)
		name := strings.TrimSuffix(rel, ".md")
		r, err := s.Get(name)
		if err != nil {
			return nil
		}
		runes = append(runes, *r)
		return nil
	})
	if err != nil {
		return nil, err
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
	if strings.ContainsAny(r.Name, " \\") {
		return fmt.Errorf("name must be a slug (no spaces or backslashes)")
	}
	for _, seg := range strings.Split(r.Name, "/") {
		if seg == "" {
			return fmt.Errorf("name contains empty segment")
		}
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	return nil
}

// write renders a rune as markdown with YAML frontmatter.
func write(path string, r Rune) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// YAML frontmatter
	fmt.Fprintln(f, "---")
	fmt.Fprintf(f, "version: %s\n", r.Version)
	fmt.Fprintf(f, "hydrated: %v\n", r.Hydrated)
	fmt.Fprintf(f, "coverage: %v\n", r.Coverage)
	fmt.Fprintln(f, "---")

	// Title and description
	fmt.Fprintf(f, "\n# %s\n\n%s\n", r.Name, r.Description)

	// Behavior
	if r.Behavior != "" {
		fmt.Fprintf(f, "\n## Behavior\n\n%s\n", r.Behavior)
	}

	// Tests
	writeList(f, "Positive tests", r.PositiveTests)
	writeList(f, "Negative tests", r.NegativeTests)

	return nil
}

func writeList(f *os.File, heading string, items []string) {
	if len(items) == 0 {
		return
	}
	fmt.Fprintf(f, "\n## %s\n\n", heading)
	for _, item := range items {
		fmt.Fprintf(f, "- %s\n", item)
	}
}

// parse reads a markdown rune file with YAML frontmatter.
func parse(content string) (*Rune, error) {
	var r Rune

	// Split frontmatter from body
	body := content
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---")
		if end != -1 {
			fm := content[4 : 4+end]
			body = content[4+end+4:] // skip past closing ---\n
			if err := yaml.Unmarshal([]byte(fm), &r); err != nil {
				return nil, fmt.Errorf("invalid frontmatter: %w", err)
			}
		}
	}

	// Parse markdown sections
	lines := strings.Split(body, "\n")
	var section string
	var sectionLines []string

	flush := func() {
		text := strings.TrimSpace(strings.Join(sectionLines, "\n"))
		switch section {
		case "":
			// Description is the first non-empty paragraph after the title
			if text != "" && r.Description == "" {
				r.Description = text
			}
		case "behavior":
			r.Behavior = text
		case "positive tests":
			r.PositiveTests = parseList(sectionLines)
		case "negative tests":
			r.NegativeTests = parseList(sectionLines)
		}
		sectionLines = nil
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") && r.Name == "" {
			r.Name = strings.TrimPrefix(line, "# ")
			section = ""
			sectionLines = nil
			continue
		}
		if strings.HasPrefix(line, "## ") {
			flush()
			section = strings.ToLower(strings.TrimPrefix(line, "## "))
			continue
		}
		sectionLines = append(sectionLines, line)
	}
	flush()

	return &r, nil
}

func parseList(lines []string) []string {
	var items []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			items = append(items, strings.TrimPrefix(line, "- "))
		}
	}
	return items
}
