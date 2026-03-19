package runner

import (
	"context"
	"fmt"

	"github.com/chrishayen/valkyrie/config"
)

// Runner executes a task inside a sandbox and returns the output.
type Runner interface {
	Run(ctx context.Context, task string) (string, error)
}

// New returns the appropriate Runner for the given agent config.
func New(agent config.Agent) (Runner, error) {
	switch agent.Type {
	case "claude-pro":
		return newClaudePro(agent), nil
	case "mock":
		return newMock(agent), nil
	case "claude-api":
		return nil, fmt.Errorf("claude-api runner not yet implemented")
	case "docker":
		return nil, fmt.Errorf("docker runner not yet implemented")
	default:
		return nil, fmt.Errorf("unknown agent type: %s", agent.Type)
	}
}
