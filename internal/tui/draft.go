package tui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Draft struct {
	ID           string           `json:"id"`
	FeatureName  string           `json:"feature_name"`
	Summary      string           `json:"summary"`
	Requirement  string           `json:"requirement"`
	Result       *decomposeResult `json:"result,omitempty"`
	Conversation []qaPair         `json:"conversation,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

type DraftStore struct {
	dir string
}

func NewDraftStore(registryPath string) *DraftStore {
	return &DraftStore{dir: filepath.Join(registryPath, ".drafts")}
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = slugRe.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "untitled"
	}
	if len(s) > 40 {
		s = s[:40]
	}
	return s
}

func (s *DraftStore) Save(d Draft) error {
	if err := os.MkdirAll(s.dir, 0755); err != nil {
		return fmt.Errorf("creating drafts dir: %w", err)
	}

	now := time.Now()
	if d.ID == "" {
		name := d.FeatureName
		if name == "" {
			name = "untitled"
		}
		d.ID = fmt.Sprintf("%s-%d", slugify(name), now.Unix())
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.UpdatedAt = now

	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling draft: %w", err)
	}
	return os.WriteFile(filepath.Join(s.dir, d.ID+".json"), data, 0644)
}

func (s *DraftStore) Load(id string) (*Draft, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, id+".json"))
	if err != nil {
		return nil, err
	}
	var d Draft
	if err := json.Unmarshal(data, &d); err != nil {
		return nil, fmt.Errorf("parsing draft: %w", err)
	}
	return &d, nil
}

func (s *DraftStore) List() ([]Draft, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var drafts []Draft
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		d, err := s.Load(id)
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

func (s *DraftStore) Delete(id string) error {
	return os.Remove(filepath.Join(s.dir, id+".json"))
}
