package runner

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chrishayen/valkyrie/config"
)

type claudeProRunner struct {
	agent config.Agent
}

func newClaudePro(agent config.Agent) *claudeProRunner {
	return &claudeProRunner{agent: agent}
}

// credentialsPath returns the resolved path to Claude credentials.
// Defaults to ~/.claude/.credentials.json if not set in config.
func (r *claudeProRunner) credentialsPath() string {
	if r.agent.CredentialsPath != "" {
		return r.agent.CredentialsPath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude", ".credentials.json")
}

// Validate checks that the claude CLI is available and credentials exist.
func (r *claudeProRunner) Validate() error {
	if _, err := exec.LookPath("claude"); err != nil {
		return fmt.Errorf("claude CLI not found in PATH — run: npm install -g @anthropic-ai/claude-code")
	}
	credPath := r.credentialsPath()
	if _, err := os.Stat(credPath); os.IsNotExist(err) {
		return fmt.Errorf("claude credentials not found at %s — run: claude setup-token", credPath)
	}
	return nil
}

// Run executes the task using the claude CLI in non-interactive mode.
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

	// Ensure ANTHROPIC_API_KEY is unset so subscription auth takes precedence
	env := os.Environ()
	filtered := env[:0]
	for _, e := range env {
		if !strings.HasPrefix(e, "ANTHROPIC_API_KEY=") {
			filtered = append(filtered, e)
		}
	}
	cmd.Env = filtered

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("claude exited with error: %w\n%s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}
