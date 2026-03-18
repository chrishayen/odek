package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chrishayen/valkyrie/config"
)

type claudeProRunner struct {
	agent config.Agent
}

func newClaudePro(agent config.Agent) *claudeProRunner {
	return &claudeProRunner{agent: agent}
}

// tokenEnv returns the name of the env var holding the OAuth token.
// Defaults to CLAUDE_CODE_OAUTH_TOKEN (produced by `claude setup-token`).
func (r *claudeProRunner) tokenEnv() string {
	if r.agent.TokenEnv != "" {
		return r.agent.TokenEnv
	}
	return "CLAUDE_CODE_OAUTH_TOKEN"
}

// Validate checks that the claude CLI is available and the token is set.
func (r *claudeProRunner) Validate() error {
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude CLI not found in PATH — run: npm install -g @anthropic-ai/claude-code")
	}
	tokenVar := r.tokenEnv()
	if os.Getenv(tokenVar) == "" {
		return fmt.Errorf("%s is not set — run: claude setup-token", tokenVar)
	}
	return nil
}

// Run executes the task using the claude CLI with OAuth subscription auth.
func (r *claudeProRunner) Run(ctx context.Context, task string) (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	args := []string{
		"--print",
		"--permission-mode", "bypassPermissions",
	}
	if r.agent.Model != "" {
		args = append(args, "--model", r.agent.Model)
	}
	args = append(args, task)

	cmd := exec.CommandContext(ctx, "claude", args...)

	// Build env from config only — no inherited env, no stripping hacks.
	cmd.Env = []string{
		"CLAUDE_CODE_OAUTH_TOKEN=" + os.Getenv(r.tokenEnv()),
		"PATH=" + os.Getenv("PATH"), // claude CLI needs PATH to find itself
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude exited with error: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
