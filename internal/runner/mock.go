package runner

import (
	"context"
	"fmt"

	"github.com/chrishayen/valkyrie/config"
)

// mockRunner is a test sandbox that returns deterministic generated code.
type mockRunner struct {
	agent config.Agent
}

func newMock(agent config.Agent) *mockRunner {
	return &mockRunner{agent: agent}
}

// Run returns a canned Go implementation + test for any task.
func (r *mockRunner) Run(_ context.Context, task string) (string, error) {
	return fmt.Sprintf(`=== FILE: go.mod ===
module valkyrie-rune

go 1.22
=== END FILE ===

=== FILE: main.go ===
package main

import "fmt"

// Generated for: %s
func Run() string {
	return "Hello, World!"
}

func main() {
	fmt.Println(Run())
}
=== END FILE ===

=== FILE: main_test.go ===
package main

import "testing"

func TestRun(t *testing.T) {
	got := Run()
	want := "Hello, World!"
	if got != want {
		t.Errorf("Run() = %%q, want %%q", got, want)
	}
}
=== END FILE ===
`, task), nil
}
