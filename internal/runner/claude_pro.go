package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/chrishayen/valkyrie/config"
)

const defaultClaudeImage = "ghcr.io/chrishayen/valkyrie-claude:latest"

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

func (r *claudeProRunner) image() string {
	if r.agent.Image != "" {
		return r.agent.Image
	}
	return defaultClaudeImage
}

// Validate checks that Docker is available and a token is resolvable.
func (r *claudeProRunner) Validate() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found in PATH — install Docker to use claude-pro sandbox")
	}
	if r.resolveToken() == "" {
		return fmt.Errorf("no token configured — set token in config or export CLAUDE_CODE_OAUTH_TOKEN (run: claude setup-token)")
	}
	return nil
}

// Run executes the task inside a Docker container with the claude CLI.
func (r *claudeProRunner) Run(ctx context.Context, task string) (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	args := []string{
		"run", "--rm",
		"-e", "CLAUDE_CODE_OAUTH_TOKEN=" + r.resolveToken(),
		r.image(),
		"--print",
		"--permission-mode", "bypassPermissions",
	}
	if r.agent.Model != "" {
		args = append(args, "--model", r.agent.Model)
	}
	args = append(args, task)

	cmd := exec.CommandContext(ctx, "docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker run failed: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
