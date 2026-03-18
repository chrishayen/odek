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

// resolveToken returns the OAuth token using this precedence:
// 1. token field in config (literal value)
// 2. token_env field in config (named env var)
// 3. CLAUDE_CODE_OAUTH_TOKEN env var (default)
func (r *claudeProRunner) resolveToken() string {
	if r.agent.Token != "" {
		return r.agent.Token
	}
	envVar := r.agent.TokenEnv
	if envVar == "" {
		envVar = "CLAUDE_CODE_OAUTH_TOKEN"
	}
	return os.Getenv(envVar)
}

// Validate checks that the claude CLI is available and a token is resolvable.
func (r *claudeProRunner) Validate() error {
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude CLI not found in PATH — run: npm install -g @anthropic-ai/claude-code")
	}
	if r.resolveToken() == "" {
		return fmt.Errorf("no token configured — set token in config or export CLAUDE_CODE_OAUTH_TOKEN (run: claude setup-token)")
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
		"CLAUDE_CODE_OAUTH_TOKEN=" + r.resolveToken(),
		"PATH=" + os.Getenv("PATH"), // claude CLI needs PATH to find itself
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude exited with error: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
