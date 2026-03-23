package runner

import (
	"context"
	"fmt"
	"io"

	"github.com/chrishayen/valkyrie/config"
)

// Runner executes a task inside a sandbox and returns the output.
// If logOut is non-nil, output is streamed to it in real time.
type Runner interface {
	Run(ctx context.Context, task string, logOut io.Writer) (string, error)
}

// New returns the appropriate Runner for the given agent config.
func New(agent config.Agent) (Runner, error) {
	switch agent.Type {
	case "claude-sub":
		return newClaudeSub(agent), nil
	case "mock":
		return newMock(agent), nil
	default:
		return nil, fmt.Errorf("unknown agent type: %s", agent.Type)
	}
}
