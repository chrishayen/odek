package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rune is a function specification identified by a fully-qualified dot-notation name.
type Rune struct {
	FQN           string   `json:"fqn"                     yaml:"-"`
	Description   string   `json:"description"              yaml:"-"`
	Signature     string   `json:"signature"                yaml:"-"`
	Behavior      string   `json:"behavior,omitempty"       yaml:"-"`
	PositiveTests []string `json:"positive_tests,omitempty" yaml:"-"`
	NegativeTests []string `json:"negative_tests,omitempty" yaml:"-"`
	Version       string   `json:"version"                  yaml:"version"`
	Project       string   `json:"project,omitempty"        yaml:"project,omitempty"`
	Status        string   `json:"status"                   yaml:"status"`
}

// RuneStore manages rune specs on the server filesystem.
type RuneStore struct {
	dataDir string // root data directory
}

func NewRuneStore(dataDir string) *RuneStore {
	return &RuneStore{dataDir: dataDir}
}

func (s *RuneStore) runesDir() string {
	return filepath.Join(s.dataDir, "runes")
}

// fqnToPath converts dot-notation FQN to filesystem path.
// "net.http.parse_url" -> "{dataDir}/runes/net/http/parse_url.md"
func (s *RuneStore) fqnToPath(fqn string) string {
	parts := strings.Split(fqn, ".")
	segments := append([]string{s.runesDir()}, parts...)
	return strings.Join(segments, string(filepath.Separator)) + ".md"
}

// pathToFQN converts a filesystem path back to dot-notation FQN.
func (s *RuneStore) pathToFQN(path string) (string, error) {
	rel, err := filepath.Rel(s.runesDir(), path)
	if err != nil {
		return "", err
	}
	rel = strings.TrimSuffix(rel, ".md")
	return strings.ReplaceAll(rel, string(filepath.Separator), "."), nil
}

func (s *RuneStore) Create(r Rune) error {
	if err := validateRune(r); err != nil {
		return err
	}
	p := s.fqnToPath(r.FQN)
	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("rune %q already exists", r.FQN)
	}
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	if r.Version == "" {
		r.Version = "0.1.0"
	}
	if r.Status == "" {
		r.Status = "draft"
	}
	return writeRune(p, r)
}

func (s *RuneStore) Get(fqn string) (*Rune, error) {
	data, err := os.ReadFile(s.fqnToPath(fqn))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("rune %q not found", fqn)
		}
		return nil, err
	}
	r, err := parseRune(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing rune %q: %w", fqn, err)
	}
	return r, nil
}

func (s *RuneStore) Update(r Rune) error {
	if err := validateRune(r); err != nil {
		return err
	}
	p := s.fqnToPath(r.FQN)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("rune %q not found", r.FQN)
	}
	return writeRune(p, r)
}

func (s *RuneStore) Delete(fqn string) error {
	p := s.fqnToPath(fqn)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("rune %q not found", fqn)
	}
	return os.Remove(p)
}

// List returns all runes, optionally filtered by project and/or namespace prefix.
func (s *RuneStore) List(project, namespace string) ([]Rune, error) {
	root := s.runesDir()
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return nil, nil
	}

	var runes []Rune
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		r, err := parseRune(string(data))
		if err != nil {
			return nil
		}

		// Apply filters
		if project != "" && r.Project != project {
			return nil
		}
		if namespace != "" && !strings.HasPrefix(r.FQN, namespace) {
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

// Search finds runes whose FQN, description, or signature contain the query string.
func (s *RuneStore) Search(query string) ([]Rune, error) {
	all, err := s.List("", "")
	if err != nil {
		return nil, err
	}
	q := strings.ToLower(query)
	var matches []Rune
	for _, r := range all {
		if strings.Contains(strings.ToLower(r.FQN), q) ||
			strings.Contains(strings.ToLower(r.Description), q) ||
			strings.Contains(strings.ToLower(r.Signature), q) {
			matches = append(matches, r)
		}
	}
	return matches, nil
}

// ListProjects returns all project names (top-level directories that contain runes with a project field).
func (s *RuneStore) ListProjects() ([]string, error) {
	projectsDir := filepath.Join(s.dataDir, "projects")
	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		return nil, nil
	}
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, err
	}
	var projects []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
			projects = append(projects, strings.TrimSuffix(e.Name(), ".toml"))
		}
	}
	return projects, nil
}

// CreateProject creates project metadata.
func (s *RuneStore) CreateProject(name string) error {
	dir := filepath.Join(s.dataDir, "projects")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, name+".toml")
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("project %q already exists", name)
	}
	content := fmt.Sprintf("name = %q\n", name)
	return os.WriteFile(path, []byte(content), 0644)
}

func validateRune(r Rune) error {
	if r.FQN == "" {
		return fmt.Errorf("fqn is required")
	}
	parts := strings.Split(r.FQN, ".")
	if len(parts) < 2 {
		return fmt.Errorf("fqn must have at least two segments (e.g. net.parse_url)")
	}
	for _, seg := range parts {
		if seg == "" {
			return fmt.Errorf("fqn contains empty segment")
		}
		if strings.ContainsAny(seg, " /\\") {
			return fmt.Errorf("fqn segments must not contain spaces or slashes")
		}
	}
	if r.Description == "" {
		return fmt.Errorf("description is required")
	}
	if r.Signature == "" {
		return fmt.Errorf("signature is required")
	}
	return nil
}

func writeRune(path string, r Rune) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "---")
	fmt.Fprintf(f, "version: %s\n", r.Version)
	if r.Project != "" {
		fmt.Fprintf(f, "project: %s\n", r.Project)
	}
	fmt.Fprintf(f, "status: %s\n", r.Status)
	fmt.Fprintln(f, "---")

	fmt.Fprintf(f, "\n# %s\n\n%s\n", r.FQN, r.Description)

	if r.Signature != "" {
		fmt.Fprintf(f, "\n## Signature\n\n%s\n", r.Signature)
	}
	if r.Behavior != "" {
		fmt.Fprintf(f, "\n## Behavior\n\n%s\n", r.Behavior)
	}
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

func parseRune(content string) (*Rune, error) {
	var r Rune

	body := content
	if strings.HasPrefix(content, "---\n") {
		end := strings.Index(content[4:], "\n---")
		if end != -1 {
			fm := content[4 : 4+end]
			body = content[4+end+4:]
			if err := yaml.Unmarshal([]byte(fm), &r); err != nil {
				return nil, fmt.Errorf("invalid frontmatter: %w", err)
			}
		}
	}

	lines := strings.Split(body, "\n")
	var section string
	var sectionLines []string

	flush := func() {
		text := strings.TrimSpace(strings.Join(sectionLines, "\n"))
		switch section {
		case "":
			if text != "" && r.Description == "" {
				r.Description = text
			}
		case "signature":
			r.Signature = text
		case "behavior":
			r.Behavior = text
		case "positive tests":
			r.PositiveTests = parseListItems(sectionLines)
		case "negative tests":
			r.NegativeTests = parseListItems(sectionLines)
		}
		sectionLines = nil
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") && r.FQN == "" {
			r.FQN = strings.TrimPrefix(line, "# ")
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

func parseListItems(lines []string) []string {
	var items []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") {
			items = append(items, strings.TrimPrefix(line, "- "))
		}
	}
	return items
}
