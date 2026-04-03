package draft

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chrishayen/odek/internal/frontmatter"
	runepkg "github.com/chrishayen/odek/internal/rune"
)

// Draft represents an in-progress feature decomposition stored as a folder
// under runes/draft/ in the registry, containing a feature.md and rune files.
type Draft struct {
	ID          string    // folder name, e.g. "hello_world_a8f3b2"
	FeatureName string    // from feature.md
	Summary     string    // from feature.md body
	Requirement string    // from feature.md body
	Version     string    // from feature.md frontmatter
	Status      string    // from feature.md frontmatter
	CreatedAt   time.Time // folder mod time
	UpdatedAt   time.Time // feature.md mod time
}

// Store manages draft folders inside runes/draft/.
type Store struct {
	dir          string // runes/draft/
	registryPath string // registry root (for merge target)
	outputPath   string // code output root
}

// NewStore creates a draft store. registryPath is the registry root.
func NewStore(registryPath, outputPath string) *Store {
	return &Store{
		dir:          filepath.Join(registryPath, "runes", "draft"),
		registryPath: registryPath,
		outputPath:   outputPath,
	}
}

// draftDir returns the path to a specific draft folder.
func (s *Store) draftDir(id string) string {
	return filepath.Join(s.dir, id)
}

// featurePath returns the path to a draft's feature.md.
func (s *Store) featurePath(id string) string {
	return filepath.Join(s.dir, id, "feature.md")
}

// runeStore returns a rune.Store scoped to a draft folder.
func (s *Store) runeStore(id string) *runepkg.Store {
	return runepkg.NewStore(s.draftDir(id), s.outputPath)
}

// Create creates a new draft folder with a feature.md.
func (s *Store) Create(featureName, requirement, summary string) (*Draft, error) {
	if featureName == "" {
		featureName = "untitled"
	}

	id := slugify(featureName) + "_" + shortID()
	dir := s.draftDir(id)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating draft dir: %w", err)
	}

	if err := s.writeFeatureMD(id, featureName, requirement, summary); err != nil {
		os.RemoveAll(dir)
		return nil, err
	}

	now := time.Now()
	return &Draft{
		ID:          id,
		FeatureName: featureName,
		Summary:     summary,
		Requirement: requirement,
		Version:     "0.1.0",
		Status:      "draft",
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update rewrites the feature.md for an existing draft.
func (s *Store) Update(id, featureName, requirement, summary string) error {
	if _, err := os.Stat(s.draftDir(id)); os.IsNotExist(err) {
		return fmt.Errorf("draft %q not found", id)
	}
	return s.writeFeatureMD(id, featureName, requirement, summary)
}

func (s *Store) writeFeatureMD(id, featureName, requirement, summary string) error {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("version: 0.1.0\n")
	b.WriteString("status: draft\n")
	b.WriteString("hydrated: false\n")
	b.WriteString("coverage: -1\n")
	b.WriteString("---\n\n")
	fmt.Fprintf(&b, "# %s\n\n", featureName)
	if requirement != "" {
		b.WriteString(requirement + "\n\n")
	}
	if summary != "" {
		b.WriteString(summary + "\n")
	}
	return os.WriteFile(s.featurePath(id), []byte(b.String()), 0644)
}

// Get reads a draft's metadata from its feature.md.
func (s *Store) Get(id string) (*Draft, error) {
	fp := s.featurePath(id)
	data, err := os.ReadFile(fp)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("draft %q not found", id)
		}
		return nil, err
	}

	var fm struct {
		Version string `yaml:"version"`
		Status  string `yaml:"status"`
	}
	body := frontmatter.Parse(string(data), &fm)

	featureName, requirement, summary := parseBody(body)

	info, _ := os.Stat(fp)
	dirInfo, _ := os.Stat(s.draftDir(id))

	d := &Draft{
		ID:          id,
		FeatureName: featureName,
		Summary:     summary,
		Requirement: requirement,
		Version:     fm.Version,
		Status:      fm.Status,
	}
	if info != nil {
		d.UpdatedAt = info.ModTime()
	}
	if dirInfo != nil {
		d.CreatedAt = dirInfo.ModTime()
	}
	return d, nil
}

// parseBody extracts the feature name (from # heading), requirement, and summary
// from the body after frontmatter.
func parseBody(body string) (featureName, requirement, summary string) {
	body = strings.TrimSpace(body)
	lines := strings.Split(body, "\n")

	// First line should be "# name"
	if len(lines) > 0 && strings.HasPrefix(lines[0], "# ") {
		featureName = strings.TrimPrefix(lines[0], "# ")
		lines = lines[1:]
	}

	// Remaining text: split into paragraphs. First is requirement, second is summary.
	remaining := strings.TrimSpace(strings.Join(lines, "\n"))
	parts := strings.SplitN(remaining, "\n\n", 2)
	if len(parts) >= 1 {
		requirement = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 {
		summary = strings.TrimSpace(parts[1])
	}
	return
}

// List returns all drafts sorted by most recently updated.
func (s *Store) List() ([]Draft, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var drafts []Draft
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		d, err := s.Get(e.Name())
		if err != nil {
			continue
		}
		drafts = append(drafts, *d)
	}

	sort.Slice(drafts, func(i, j int) bool {
		return drafts[i].UpdatedAt.After(drafts[j].UpdatedAt)
	})
	return drafts, nil
}

// Delete removes a draft folder entirely.
func (s *Store) Delete(id string) error {
	dir := s.draftDir(id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("draft %q not found", id)
	}
	return os.RemoveAll(dir)
}

// SaveRunes writes rune files into the draft's tree structure.
// Existing runes in the draft are cleared first.
func (s *Store) SaveRunes(id string, runes []runepkg.Rune) error {
	dir := s.draftDir(id)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("draft %q not found", id)
	}

	// Clear existing rune dirs (preserve feature.md)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() {
			os.RemoveAll(filepath.Join(dir, e.Name()))
		}
	}

	rs := s.runeStore(id)
	for _, r := range runes {
		if err := rs.Create(r); err != nil {
			return fmt.Errorf("saving rune %q: %w", r.Name, err)
		}
	}
	return nil
}

// ListRunes returns all runes in a draft folder.
func (s *Store) ListRunes(id string) ([]runepkg.Rune, error) {
	rs := s.runeStore(id)
	return rs.List()
}

// Merge moves a draft's runes into the main registry.
// The draft folder is deleted after a successful merge.
func (s *Store) Merge(id string) error {
	_, err := s.Get(id)
	if err != nil {
		return err
	}

	runes, err := s.ListRunes(id)
	if err != nil {
		return fmt.Errorf("listing draft runes: %w", err)
	}

	// Create runes in the main registry
	mainRuneStore := runepkg.NewStore(filepath.Join(s.registryPath, "runes"), s.outputPath)
	for _, r := range runes {
		if err := mainRuneStore.Create(r); err != nil {
			if strings.Contains(err.Error(), "already exists") {
				continue
			}
			return fmt.Errorf("merging rune %q: %w", r.Name, err)
		}
	}

	// Delete the draft folder
	return s.Delete(id)
}

// --- helpers ---

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRe.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	if s == "" {
		return "untitled"
	}
	if len(s) > 40 {
		s = s[:40]
	}
	return s
}

func shortID() string {
	b := make([]byte, 3)
	rand.Read(b)
	return hex.EncodeToString(b)
}
