package runner

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/chrishayen/valkyrie/config"
)

const defaultClaudeImage = "ghcr.io/chrishayen/valkyrie-claude:latest"

type claudeSubRunner struct {
	agent config.Agent
}

func newClaudeSub(agent config.Agent) *claudeSubRunner {
	return &claudeSubRunner{agent: agent}
}

// resolveToken returns the OAuth token using this precedence:
// 1. token field in config (literal value)
// 2. token_env field in config (named env var)
// 3. CLAUDE_CODE_OAUTH_TOKEN env var
// 4. Parse from ~/.claude/.credentials.json
func (r *claudeSubRunner) resolveToken() string {
	if r.agent.Token != "" {
		return r.agent.Token
	}
	if r.agent.TokenEnv != "" {
		if v := os.Getenv(r.agent.TokenEnv); v != "" {
			return v
		}
	}
	if v := os.Getenv("CLAUDE_CODE_OAUTH_TOKEN"); v != "" {
		return v
	}
	return config.ClaudeToken()
}

func (r *claudeSubRunner) image() string {
	if r.agent.Image != "" {
		return r.agent.Image
	}
	return defaultClaudeImage
}

// Validate checks that Docker is available and a token is resolvable.
func (r *claudeSubRunner) Validate() error {
	if _, err := exec.LookPath("docker"); err != nil {
		return fmt.Errorf("docker not found in PATH — install Docker to use claude-sub")
	}
	if r.resolveToken() == "" {
		return fmt.Errorf("no token found — set token in config, export CLAUDE_CODE_OAUTH_TOKEN, or log in with 'claude login'")
	}
	return nil
}

// Run executes the task inside a Docker container with the claude CLI.
// If logOut is non-nil, the full stream-json event stream is written to it in real time.
// The return value contains only the final result text extracted from the stream.
func (r *claudeSubRunner) Run(ctx context.Context, task string, logOut io.Writer) (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}

	args := []string{
		"run", "--rm",
		"-e", "CLAUDE_CODE_OAUTH_TOKEN=" + r.resolveToken(),
		r.image(),
		"--output-format", "stream-json",
		"--verbose",
		"--permission-mode", "bypassPermissions",
	}
	if r.agent.Model != "" {
		args = append(args, "--model", r.agent.Model)
	}
	args = append(args, task)

	cmd := exec.CommandContext(ctx, "docker", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("creating stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout // merge stderr into stdout

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("starting docker: %w", err)
	}

	var allLines bytes.Buffer
	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 10MB line buffer
	for scanner.Scan() {
		line := scanner.Bytes()
		allLines.Write(line)
		allLines.WriteByte('\n')
		if logOut != nil {
			formatEvent(logOut, line)
		}
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("docker run failed: %w\n%s", err, allLines.String())
	}

	return extractResult(allLines.Bytes())
}

// streamEvent captures the fields we care about from stream-json events.
type streamEvent struct {
	Type    string `json:"type"`
	Subtype string `json:"subtype,omitempty"`
	Result  string `json:"result,omitempty"`

	// assistant message fields
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`
		Name string `json:"name,omitempty"` // tool name
	} `json:"content,omitempty"`

	// tool_use / tool_result
	Name  string `json:"name,omitempty"`
	Input json.RawMessage `json:"input,omitempty"`
}

// extractResult parses stream-json output and returns the text from the final result event.
func extractResult(data []byte) (string, error) {
	var resultText string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024)
	for scanner.Scan() {
		var ev streamEvent
		if json.Unmarshal(scanner.Bytes(), &ev) != nil {
			continue
		}
		if ev.Type == "result" {
			resultText = ev.Result
		}
	}

	if resultText == "" {
		return "", fmt.Errorf("no result event found in stream output")
	}
	return resultText, nil
}

// formatEvent writes a human-readable summary of a stream-json event to w.
func formatEvent(w io.Writer, line []byte) {
	var ev streamEvent
	if json.Unmarshal(line, &ev) != nil {
		return
	}

	switch ev.Type {
	case "system":
		if ev.Subtype == "init" {
			fmt.Fprintf(w, "[agent] session started\n")
		}
	case "assistant":
		for _, c := range ev.Content {
			switch c.Type {
			case "text":
				if c.Text != "" {
					fmt.Fprintf(w, "[assistant] %s\n", c.Text)
				}
			case "tool_use":
				fmt.Fprintf(w, "[tool] %s\n", c.Name)
			}
		}
	case "tool_use":
		fmt.Fprintf(w, "[tool] %s\n", ev.Name)
	case "tool_result":
		// skip — too noisy
	case "result":
		fmt.Fprintf(w, "[agent] done\n")
	}
}
