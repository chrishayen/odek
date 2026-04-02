package jobs

import (
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// JobStatus represents the current state of a job.
type JobStatus string

const (
	StatusPending   JobStatus = "pending"
	StatusRunning   JobStatus = "running"
	StatusCompleted JobStatus = "completed"
	StatusFailed    JobStatus = "failed"
)

// Job represents an async task with status tracking.
type Job struct {
	ID        string          `json:"id"`
	Status    JobStatus       `json:"status"`
	Progress  string          `json:"progress,omitempty"`
	Result    json.RawMessage `json:"result,omitempty"`
	Error     string          `json:"error,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// Manager tracks async jobs.
type Manager struct {
	jobs    sync.Map
	counter atomic.Int64
}

// Create creates a new pending job and returns it.
func (m *Manager) Create() *Job {
	id := fmt.Sprintf("job_%d", m.counter.Add(1))
	j := &Job{
		ID:        id,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}
	m.jobs.Store(id, j)
	return j
}

// Get returns a job by ID.
func (m *Manager) Get(id string) (*Job, bool) {
	v, ok := m.jobs.Load(id)
	if !ok {
		return nil, false
	}
	return v.(*Job), true
}

// SetProgress updates a job's progress message.
func (m *Manager) SetProgress(id string, msg string) {
	if v, ok := m.jobs.Load(id); ok {
		j := v.(*Job)
		j.Progress = msg
	}
}

// SetRunning updates a job's status to running.
func (m *Manager) SetRunning(id string) {
	if v, ok := m.jobs.Load(id); ok {
		j := v.(*Job)
		j.Status = StatusRunning
	}
}

// SetCompleted marks a job as completed with a result.
func (m *Manager) SetCompleted(id string, result json.RawMessage) {
	if v, ok := m.jobs.Load(id); ok {
		j := v.(*Job)
		j.Status = StatusCompleted
		j.Result = result
	}
}

// SetFailed marks a job as failed with an error message.
func (m *Manager) SetFailed(id string, err error) {
	if v, ok := m.jobs.Load(id); ok {
		j := v.(*Job)
		j.Status = StatusFailed
		j.Error = err.Error()
	}
}
