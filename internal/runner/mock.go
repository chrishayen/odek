package runner

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/chrishayen/valkyrie/config"
)

// mockRunner is a test sandbox that returns deterministic generated code.
type mockRunner struct {
	agent config.Agent
}

func newMock(agent config.Agent) *mockRunner {
	return &mockRunner{agent: agent}
}

// Run returns a canned response. If the task looks like an analyze prompt,
// it returns a JSON analysis result. Otherwise it returns generated code.
// If logOut is non-nil, the response is also streamed there.
func (r *mockRunner) Run(_ context.Context, task string, logOut io.Writer) (string, error) {
	if strings.Contains(task, "Requirements to analyze") {
		out := mockAnalyzeResponse()
		if logOut != nil {
			io.WriteString(logOut, out)
		}
		return out, nil
	}
	hydrationOut := fmt.Sprintf(`=== FILE: go.mod ===
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
`, task)
	if logOut != nil {
		io.WriteString(logOut, hydrationOut)
	}
	return hydrationOut, nil
}

func mockAnalyzeResponse() string {
	return `{
  "new_runes": [
    {
      "name": "auth/validate-email",
      "description": "Validates that a string is a well-formed email address.",
      "behavior": "Accepts a string. Returns true if the string contains exactly one @ symbol with non-empty local and domain parts. Returns false otherwise.",
      "positive_tests": [
        "Given 'user@example.com', returns true",
        "Given 'a@b.co', returns true"
      ],
      "negative_tests": [
        "Given an empty string, returns false",
        "Given 'no-at-sign', returns false",
        "Given '@missing-local.com', returns false"
      ]
    },
    {
      "name": "auth/hash-password",
      "description": "Produces a one-way hash of a plaintext password string.",
      "behavior": "Accepts a non-empty string. Returns a deterministic hash string. Same input always produces the same output.",
      "positive_tests": [
        "Given 'secret123', returns a non-empty string",
        "Given the same input twice, returns the same hash both times"
      ],
      "negative_tests": [
        "Given an empty string, returns an error"
      ]
    },
    {
      "name": "auth/store-credentials",
      "description": "Persists user credentials to the database.",
      "behavior": "Accepts a user ID and hashed password. Writes them to the credentials table. Returns success or error.",
      "positive_tests": [
        "Given valid user ID and hash, returns success"
      ],
      "negative_tests": [
        "Given empty user ID, returns an error"
      ]
    }
  ],
  "existing_runes": []
}`
}
