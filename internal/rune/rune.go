package rune

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Semver represents a major.minor.patch version.
type Semver struct {
	Major, Minor, Patch int
}

func (v Semver) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v Semver) IsZero() bool {
	return v.Major == 0 && v.Minor == 0 && v.Patch == 0
}

func ParseSemver(s string) Semver {
	var v Semver
	fmt.Sscanf(s, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	return v
}

func (v Semver) BumpMajor() Semver { return Semver{v.Major + 1, 0, 0} }
func (v Semver) BumpMinor() Semver { return Semver{v.Major, v.Minor + 1, 0} }
func (v Semver) BumpPatch() Semver { return Semver{v.Major, v.Minor, v.Patch + 1} }

func (v Semver) Less(other Semver) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}
	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}
	return v.Patch < other.Patch
}

func (v Semver) MarshalJSON() ([]byte, error)    { return json.Marshal(v.String()) }
func (v Semver) MarshalYAML() (interface{}, error) { return v.String(), nil }

func (v *Semver) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*v = ParseSemver(s)
	return nil
}

func (v *Semver) UnmarshalYAML(value *yaml.Node) error {
	*v = ParseSemver(value.Value)
	return nil
}

// IsSemverFilename checks if a filename is a semver-named file (e.g. "1.0.0.md").
func IsSemverFilename(name string) (Semver, bool) {
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	var v Semver
	n, _ := fmt.Sscanf(base, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	if n != 3 {
		return Semver{}, false
	}
	if base != v.String() {
		return Semver{}, false
	}
	return v, true
}

// LatestVersion scans a directory and returns the highest semver among .md files.
func LatestVersion(dir string) (Semver, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Semver{}, false
	}
	found := false
	var best Semver
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if v, ok := IsSemverFilename(e.Name()); ok {
			if !found || best.Less(v) {
				best = v
				found = true
			}
		}
	}
	return best, found
}

// Node represents a parsed item from a composition tree output.
type Node struct {
	Path        string   // dot path, e.g. "std.cli.parse_flags"
	Signature   string   // e.g. "(argv: list[string]) -> result[ParseFlagsResult, string]"
	Pos         []string // positive test cases
	Neg         []string // negative test cases
	Refs        []string // -> references
	Extend      bool     // true if this is a ~> extension
	Assumptions []string // ? assumptions
}

// IsDotPath checks if a string looks like a dot-notation path (e.g. "std.cli.parse_flags").
func IsDotPath(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c != '.' && c != '_' && !(c >= 'a' && c <= 'z') && !(c >= '0' && c <= '9') {
			return false
		}
	}
	return !strings.HasPrefix(s, ".") && !strings.HasSuffix(s, ".")
}

// ParseTree parses the indented composition tree output into nodes.
func ParseTree(output string) []Node {
	var nodes []Node
	var current *Node

	for _, line := range strings.Split(output, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "@ ") {
			if current != nil {
				current.Signature = strings.TrimPrefix(trimmed, "@ ")
			}
		} else if strings.HasPrefix(trimmed, "+ ") {
			if current != nil {
				current.Pos = append(current.Pos, strings.TrimPrefix(trimmed, "+ "))
			}
		} else if strings.HasPrefix(trimmed, "- ") {
			if current != nil {
				current.Neg = append(current.Neg, strings.TrimPrefix(trimmed, "- "))
			}
		} else if strings.HasPrefix(trimmed, "? ") {
			if current != nil {
				current.Assumptions = append(current.Assumptions, strings.TrimPrefix(trimmed, "? "))
			}
		} else if strings.HasPrefix(trimmed, "-> ") {
			if current != nil {
				current.Refs = append(current.Refs, strings.TrimPrefix(trimmed, "-> "))
			}
		} else if strings.HasPrefix(trimmed, "~> ") {
			extPath := strings.TrimPrefix(trimmed, "~> ")
			if IsDotPath(extPath) {
				nodes = append(nodes, Node{Path: extPath, Extend: true})
				current = &nodes[len(nodes)-1]
			}
		} else if IsDotPath(trimmed) {
			nodes = append(nodes, Node{Path: trimmed})
			current = &nodes[len(nodes)-1]
		}
	}

	return nodes
}

// BuildChildrenMap returns parent -> []child dot path mappings.
func BuildChildrenMap(dotPaths []string) map[string][]string {
	children := make(map[string][]string)
	for _, p := range dotPaths {
		parts := strings.Split(p, ".")
		if len(parts) > 1 {
			parent := strings.Join(parts[:len(parts)-1], ".")
			children[parent] = append(children[parent], p)
		}
	}
	return children
}

// ParseRef extracts path and pinned major from "std.cli.parse_flags@1".
func ParseRef(ref string) (string, int) {
	parts := strings.SplitN(ref, "@", 2)
	if len(parts) != 2 {
		return "", 0
	}
	var major int
	fmt.Sscanf(parts[1], "%d", &major)
	return parts[0], major
}

// Rune is the atomic unit of functionality.
type Rune struct {
	Name          string   `json:"name"                     yaml:"-"`
	Description   string   `json:"description"              yaml:"-"`
	Signature     string   `json:"signature"                yaml:"signature"`
	Behavior      string   `json:"behavior,omitempty"       yaml:"-"`
	PositiveTests []string `json:"positive_tests,omitempty" yaml:"-"`
	NegativeTests []string `json:"negative_tests,omitempty" yaml:"-"`
	Assumptions   []string `json:"assumptions,omitempty"    yaml:"-"`
	Version       Semver   `json:"version"                  yaml:"version"`
	Status        string   `json:"status,omitempty"         yaml:"status"`
	Hydrated      bool     `json:"hydrated"                 yaml:"hydrated"`
	Coverage      float64  `json:"coverage"                 yaml:"coverage"`
	Dependencies  []string `json:"dependencies,omitempty"   yaml:"dependencies"`
}

// Store manages rune specs on disk using dot-path directories with semver-named files.
type Store struct {
	runesPath  string
	outputPath string
}

func NewStore(runesPath, outputPath string) *Store {
	return &Store{runesPath: runesPath, outputPath: outputPath}
}

// OutputPath returns the root output directory for generated code.
func (s *Store) OutputPath() string { return s.outputPath }

// runeDir returns the directory for a dot-path rune name.
func (s *Store) runeDir(name string) string {
	return filepath.Join(s.runesPath, strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

// CodeDir returns the directory where generated code for a rune is stored.
func (s *Store) CodeDir(name string) string {
	return filepath.Join(s.outputPath, strings.ReplaceAll(name, ".", string(filepath.Separator)))
}

func (s *Store) Create(r Rune) error {
	if err := validate(r); err != nil {
		return err
	}
	dir := s.runeDir(r.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating rune dir: %w", err)
	}
	if _, found := LatestVersion(dir); found {
		return fmt.Errorf("rune %q already exists", r.Name)
	}
	if r.Version.IsZero() {
		r.Version = Semver{1, 0, 0}
	}
	return write(filepath.Join(dir, r.Version.String()+".md"), r)
}

func (s *Store) Get(name string) (*Rune, error) {
	dir := s.runeDir(name)
	ver, found := LatestVersion(dir)
	if !found {
		return nil, fmt.Errorf("rune %q not found", name)
	}
	data, err := os.ReadFile(filepath.Join(dir, ver.String()+".md"))
	if err != nil {
		return nil, err
	}
	r, err := parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing rune %q: %w", name, err)
	}
	if r.Name == "" {
		r.Name = name
	}
	return r, nil
}

func (s *Store) List() ([]Rune, error) {
	all, err := s.ScanAll()
	if err != nil {
		return nil, err
	}
	latest := make(map[string]*Rune)
	for i := range all {
		r := &all[i]
		if cur, ok := latest[r.Name]; !ok || cur.Version.Less(r.Version) {
			latest[r.Name] = r
		}
	}
	var runes []Rune
	for _, r := range latest {
		runes = append(runes, *r)
	}
	return runes, nil
}

// TopLevelPackages returns the latest rune for each top-level directory under runes/.
// These represent "features" — packages with no dot in their name that have
// a signature (not empty structural packages). Excludes drafts.
func (s *Store) TopLevelPackages() ([]Rune, error) {
	runes, err := s.List()
	if err != nil {
		return nil, err
	}
	var pkgs []Rune
	for _, r := range runes {
		if !strings.Contains(r.Name, ".") && r.Signature != "" && r.Status != "draft" {
			pkgs = append(pkgs, r)
		}
	}
	return pkgs, nil
}

// TopLevelDrafts returns top-level packages where status == "draft".
func (s *Store) TopLevelDrafts() ([]Rune, error) {
	runes, err := s.List()
	if err != nil {
		return nil, err
	}
	var drafts []Rune
	for _, r := range runes {
		if !strings.Contains(r.Name, ".") && r.Status == "draft" {
			drafts = append(drafts, r)
		}
	}
	return drafts, nil
}

// ListByStatus returns all runes matching the given status.
func (s *Store) ListByStatus(status string) ([]Rune, error) {
	runes, err := s.List()
	if err != nil {
		return nil, err
	}
	var filtered []Rune
	for _, r := range runes {
		if r.Status == status {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

// SetStatus reads a rune, changes its status, and writes it back.
func (s *Store) SetStatus(name, status string) error {
	r, err := s.Get(name)
	if err != nil {
		return err
	}
	r.Status = status
	return s.Update(*r)
}

// ListByPrefix returns all runes whose name starts with the given prefix.
func (s *Store) ListByPrefix(prefix string) ([]Rune, error) {
	runes, err := s.List()
	if err != nil {
		return nil, err
	}
	var filtered []Rune
	for _, r := range runes {
		if strings.HasPrefix(r.Name, prefix) {
			filtered = append(filtered, r)
		}
	}
	return filtered, nil
}

// DeleteByPrefix deletes all runes whose name starts with the given prefix,
// plus the rune with exactly that name.
func (s *Store) DeleteByPrefix(prefix string) error {
	runes, err := s.List()
	if err != nil {
		return err
	}
	for _, r := range runes {
		if r.Name == prefix || strings.HasPrefix(r.Name, prefix+".") {
			_ = s.Delete(r.Name)
		}
	}
	return nil
}

func (s *Store) Update(r Rune) error {
	if err := validate(r); err != nil {
		return err
	}
	dir := s.runeDir(r.Name)
	if _, found := LatestVersion(dir); !found {
		return fmt.Errorf("rune %q not found", r.Name)
	}
	return write(filepath.Join(dir, r.Version.String()+".md"), r)
}

func (s *Store) Delete(name string) error {
	dir := s.runeDir(name)
	if _, found := LatestVersion(dir); !found {
		return fmt.Errorf("rune %q not found", name)
	}
	return os.RemoveAll(dir)
}

// ScanAll walks the registry and returns ALL rune versions found.
func (s *Store) ScanAll() ([]Rune, error) {
	base := s.runesPath
	if _, err := os.Stat(base); os.IsNotExist(err) {
		return nil, nil
	}
	var runes []Rune
	err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}
		if _, ok := IsSemverFilename(d.Name()); !ok {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		r, err := parse(string(data))
		if err != nil {
			return nil
		}
		if r.Name == "" {
			rel, _ := filepath.Rel(base, filepath.Dir(path))
			r.Name = strings.ReplaceAll(rel, string(filepath.Separator), ".")
		}
		runes = append(runes, *r)
		return nil
	})
	return runes, err
}


// CheckStaleRefs scans runes for references whose pinned major is behind the current major.
func (s *Store) CheckStaleRefs() (stale, ok int, err error) {
	runes, err := s.ScanAll()
	if err != nil {
		return 0, 0, err
	}
	versions := make(map[string]Semver)
	for _, r := range runes {
		if cur, exists := versions[r.Name]; !exists || cur.Less(r.Version) {
			versions[r.Name] = r.Version
		}
	}
	for _, r := range runes {
		for _, ref := range r.Dependencies {
			refPath, pinnedMajor := ParseRef(ref)
			if refPath == "" {
				continue
			}
			if cur, exists := versions[refPath]; exists {
				if cur.Major > pinnedMajor {
					stale++
				} else {
					ok++
				}
			}
		}
	}
	return stale, ok, nil
}

// FormatExistingContext builds a prompt context string from existing runes.
func (s *Store) FormatExistingContext() (string, error) {
	runes, err := s.List()
	if err != nil {
		return "", err
	}
	if len(runes) == 0 {
		return "", nil
	}

	var sb strings.Builder
	sb.WriteString("The following units already exist in the registry. Reference existing ones with -> path@MAJOR instead of recreating them. If a new requirement needs to EXTEND an existing unit with additional capabilities, emit it as ~> path.to.unit and include only the NEW test cases to add.\n\n")

	for _, r := range runes {
		sb.WriteString(fmt.Sprintf("%s @%d (v%s)\n", r.Name, r.Version.Major, r.Version))
		if r.Signature != "" {
			sb.WriteString("  @ " + r.Signature + "\n")
		}
		for _, t := range r.PositiveTests {
			sb.WriteString("  + " + t + "\n")
		}
		for _, t := range r.NegativeTests {
			sb.WriteString("  - " + t + "\n")
		}
	}

	return sb.String(), nil
}

// --- internal ---

func validate(r Rune) error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !IsDotPath(r.Name) {
		return fmt.Errorf("name must be a dot-separated path (e.g. auth.validate_email)")
	}
	return nil
}

func write(path string, r Rune) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "---")
	fmt.Fprintf(f, "version: %s\n", r.Version)
	if r.Status != "" {
		fmt.Fprintf(f, "status: %s\n", r.Status)
	}
	fmt.Fprintf(f, "hydrated: %v\n", r.Hydrated)
	fmt.Fprintf(f, "coverage: %v\n", r.Coverage)
	if r.Signature != "" {
		fmt.Fprintf(f, "signature: '%s'\n", r.Signature)
	}
	if len(r.Dependencies) > 0 {
		fmt.Fprintln(f, "dependencies:")
		for _, d := range r.Dependencies {
			fmt.Fprintf(f, "  - %s\n", d)
		}
	}
	fmt.Fprintln(f, "---")

	fmt.Fprintf(f, "\n# %s\n\n%s\n", r.Name, r.Description)

	if r.Signature != "" {
		fmt.Fprintf(f, "\n## Signature\n\n%s\n", r.Signature)
	}

	if r.Behavior != "" {
		fmt.Fprintf(f, "\n## Behavior\n\n%s\n", r.Behavior)
	}

	if len(r.PositiveTests) > 0 || len(r.NegativeTests) > 0 {
		fmt.Fprintln(f, "\n## Tests")
		fmt.Fprintln(f)
		for _, t := range r.PositiveTests {
			fmt.Fprintf(f, "+ %s\n", t)
		}
		for _, t := range r.NegativeTests {
			fmt.Fprintf(f, "- %s\n", t)
		}
	}

	if len(r.Assumptions) > 0 {
		fmt.Fprintln(f, "\n## Assumptions")
		fmt.Fprintln(f)
		for _, a := range r.Assumptions {
			fmt.Fprintf(f, "? %s\n", a)
		}
	}

	return nil
}

func parse(content string) (*Rune, error) {
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
			if r.Signature == "" {
				r.Signature = text
			}
		case "behavior":
			r.Behavior = text
		case "tests":
			for _, line := range sectionLines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "+ ") {
					r.PositiveTests = append(r.PositiveTests, strings.TrimPrefix(trimmed, "+ "))
				} else if strings.HasPrefix(trimmed, "- ") {
					r.NegativeTests = append(r.NegativeTests, strings.TrimPrefix(trimmed, "- "))
				}
			}
		case "positive tests":
			r.PositiveTests = parseList(sectionLines)
		case "negative tests":
			r.NegativeTests = parseList(sectionLines)
		case "assumptions":
			for _, line := range sectionLines {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "? ") {
					r.Assumptions = append(r.Assumptions, strings.TrimPrefix(trimmed, "? "))
				} else if trimmed != "" {
					r.Assumptions = append(r.Assumptions, trimmed)
				}
			}
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

