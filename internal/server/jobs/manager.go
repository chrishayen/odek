package jobs

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Job tracks an async requirements processing job.
type Job struct {
	ID        string          `json:"id"`
	Status    Status          `json:"status"`
	Error     string          `json:"error,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Manager tracks active jobs.
type Manager struct {
	jobs sync.Map
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Create() *Job {
	j := &Job{
		ID:        newID(),
		Status:    StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.jobs.Store(j.ID, j)
	return j
}

func (m *Manager) Get(id string) *Job {
	v, ok := m.jobs.Load(id)
	if !ok {
		return nil
	}
	return v.(*Job)
}

func (m *Manager) SetRunning(id string) {
	if j := m.Get(id); j != nil {
		j.Status = StatusRunning
		j.UpdatedAt = time.Now()
	}
}

func (m *Manager) SetCompleted(id string, result json.RawMessage) {
	if j := m.Get(id); j != nil {
		j.Status = StatusCompleted
		j.Result = result
		j.UpdatedAt = time.Now()
	}
}

func (m *Manager) SetFailed(id string, err error) {
	if j := m.Get(id); j != nil {
		j.Status = StatusFailed
		j.Error = err.Error()
		j.UpdatedAt = time.Now()
	}
}

func newID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
