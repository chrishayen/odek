// Package toollog writes JSONL entries for find_examples tool calls so we can
// monitor which examples the decompose agent is actually retrieving and
// debug bad decompositions after the fact.
package toollog

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Entry is one line of the log — a single find_examples invocation.
type Entry struct {
	TS          string   `json:"ts"`
	Requirement string   `json:"requirement"`
	Query       string   `json:"query"`
	MaxResults  int      `json:"max_results"`
	Results     []string `json:"results"`
}

// Logger writes Entry values to an append-only JSONL file. Safe for
// concurrent use — the expansion phase dispatches parallel decompose calls
// that may log at the same time.
type Logger struct {
	f  *os.File
	mu sync.Mutex
}

// NewLogger opens path in append mode, creating it if missing.
func NewLogger(path string) (*Logger, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	return &Logger{f: f}, nil
}

// LogToolCall appends one JSONL entry for a find_examples call. ts may be the
// zero Time to use time.Now().
func (l *Logger) LogToolCall(ts time.Time, requirement, query string, maxResults int, paths []string) error {
	if ts.IsZero() {
		ts = time.Now()
	}
	entry := Entry{
		TS:          ts.UTC().Format(time.RFC3339),
		Requirement: requirement,
		Query:       query,
		MaxResults:  maxResults,
		Results:     paths,
	}
	buf, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	buf = append(buf, '\n')

	l.mu.Lock()
	defer l.mu.Unlock()
	if _, err := l.f.Write(buf); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

// Close flushes and closes the underlying file.
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.f.Close()
}
