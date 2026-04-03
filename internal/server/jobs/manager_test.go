package jobs

import (
	"encoding/json"
	"testing"
)

func TestManagerCreateAndGet(t *testing.T) {
	m := &Manager{}

	j := m.Create()
	if j.ID == "" {
		t.Error("expected non-empty ID")
	}
	if j.Status != StatusPending {
		t.Errorf("Status = %q, want pending", j.Status)
	}

	got, ok := m.Get(j.ID)
	if !ok {
		t.Fatal("expected to find job")
	}
	if got.ID != j.ID {
		t.Errorf("ID mismatch")
	}
}

func TestManagerGetNotFound(t *testing.T) {
	m := &Manager{}
	_, ok := m.Get("nonexistent")
	if ok {
		t.Error("expected not found")
	}
}

func TestManagerLifecycle(t *testing.T) {
	m := &Manager{}
	j := m.Create()

	m.SetRunning(j.ID)
	got, _ := m.Get(j.ID)
	if got.Status != StatusRunning {
		t.Errorf("after SetRunning: Status = %q", got.Status)
	}

	m.SetProgress(j.ID, "50% done")
	got, _ = m.Get(j.ID)
	if got.Progress != "50% done" {
		t.Errorf("Progress = %q", got.Progress)
	}

	result := json.RawMessage(`{"count": 5}`)
	m.SetCompleted(j.ID, result)
	got, _ = m.Get(j.ID)
	if got.Status != StatusCompleted {
		t.Errorf("after SetCompleted: Status = %q", got.Status)
	}
	if string(got.Result) != `{"count": 5}` {
		t.Errorf("Result = %s", got.Result)
	}
}

func TestManagerSetFailed(t *testing.T) {
	m := &Manager{}
	j := m.Create()

	m.SetRunning(j.ID)
	m.SetFailed(j.ID, &testErr{msg: "something broke"})

	got, _ := m.Get(j.ID)
	if got.Status != StatusFailed {
		t.Errorf("Status = %q", got.Status)
	}
	if got.Error != "something broke" {
		t.Errorf("Error = %q", got.Error)
	}
}

func TestManagerUniqueIDs(t *testing.T) {
	m := &Manager{}
	j1 := m.Create()
	j2 := m.Create()
	if j1.ID == j2.ID {
		t.Error("expected unique IDs")
	}
}

type testErr struct{ msg string }

func (e *testErr) Error() string { return e.msg }
