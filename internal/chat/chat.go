package chat

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Role is a message sender.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is a single turn in the conversation.
type Message struct {
	Role      Role      `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// Signal is an action the chat wants the parent to perform.
type Signal struct {
	Type string // "refine_feature", "refine_rune", "apply_changes"
	Data string // refinement text or other payload
}

// Context describes what the chat is about.
type Context struct {
	FeatureName    string `json:"feature_name,omitempty"`
	FeatureSummary string `json:"feature_summary,omitempty"`
	RuneName       string `json:"rune_name,omitempty"`
	RuneSignature  string `json:"rune_signature,omitempty"`
	RuneDesc       string `json:"rune_description,omitempty"`
	Requirement    string `json:"requirement,omitempty"`
	TreeOutput     string `json:"tree_output,omitempty"`
}

// FormatSystem renders the context as a system-level description for the LLM.
func (c Context) FormatSystem() string {
	var b strings.Builder
	if c.FeatureName != "" {
		b.WriteString("Feature: " + c.FeatureName + "\n")
	}
	if c.FeatureSummary != "" {
		b.WriteString("Summary: " + c.FeatureSummary + "\n")
	}
	if c.Requirement != "" {
		b.WriteString("Requirement: " + c.Requirement + "\n")
	}
	if c.RuneName != "" {
		b.WriteString("\nRune: " + c.RuneName + "\n")
		if c.RuneSignature != "" {
			b.WriteString("Signature: " + c.RuneSignature + "\n")
		}
		if c.RuneDesc != "" {
			b.WriteString("Description: " + c.RuneDesc + "\n")
		}
	}
	if c.TreeOutput != "" {
		b.WriteString("\nDecomposition:\n" + c.TreeOutput + "\n")
	}
	return b.String()
}

// Session is a persistent conversation.
type Session struct {
	ID        string    `json:"id"`
	Context   Context   `json:"context"`
	Messages  []Message `json:"messages"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddUser appends a user message.
func (s *Session) AddUser(content string) {
	s.Messages = append(s.Messages, Message{
		Role:      RoleUser,
		Content:   content,
		Timestamp: time.Now(),
	})
	s.UpdatedAt = time.Now()
}

// AddAssistant appends an assistant message.
func (s *Session) AddAssistant(content string) {
	s.Messages = append(s.Messages, Message{
		Role:      RoleAssistant,
		Content:   content,
		Timestamp: time.Now(),
	})
	s.UpdatedAt = time.Now()
}

// Store manages chat sessions on disk.
type Store struct {
	dir string // e.g. .odek/chats
}

// NewStore creates a Store rooted at registryPath/.odek/chats.
func NewStore(registryPath string) *Store {
	dir := filepath.Join(registryPath, ".odek", "chats")
	return &Store{dir: dir}
}

func (st *Store) ensureDir() error {
	return os.MkdirAll(st.dir, 0o755)
}

func randomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// New creates a new session with the given context.
func (st *Store) New(ctx Context) *Session {
	now := time.Now()
	return &Session{
		ID:        randomID(),
		Context:   ctx,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Save writes a session to disk.
func (st *Store) Save(s *Session) error {
	if err := st.ensureDir(); err != nil {
		return fmt.Errorf("chat store dir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}
	path := filepath.Join(st.dir, s.ID+".json")
	return os.WriteFile(path, data, 0o644)
}

// Load reads a session by ID.
func (st *Store) Load(id string) (*Session, error) {
	path := filepath.Join(st.dir, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read session %s: %w", id, err)
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}
	return &s, nil
}

// FindByFeature returns sessions that match the given feature name, newest first.
func (st *Store) FindByFeature(featureName string) ([]*Session, error) {
	return st.findBy(func(s *Session) bool {
		return s.Context.FeatureName == featureName
	})
}

// FindByRune returns sessions for a specific rune within a feature.
func (st *Store) FindByRune(featureName, runeName string) ([]*Session, error) {
	return st.findBy(func(s *Session) bool {
		return s.Context.FeatureName == featureName && s.Context.RuneName == runeName
	})
}

func (st *Store) findBy(match func(*Session) bool) ([]*Session, error) {
	if err := st.ensureDir(); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(st.dir)
	if err != nil {
		return nil, err
	}
	var sessions []*Session
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(e.Name(), ".json")
		s, err := st.Load(id)
		if err != nil {
			continue
		}
		if match(s) {
			sessions = append(sessions, s)
		}
	}
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})
	return sessions, nil
}

// Delete removes a session file.
func (st *Store) Delete(id string) error {
	return os.Remove(filepath.Join(st.dir, id+".json"))
}
