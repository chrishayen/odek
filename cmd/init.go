package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Valkyrie project in the current directory",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		tomlPath := filepath.Join(cwd, "valkyrie.toml")
		if _, err := os.Stat(tomlPath); err == nil {
			return fmt.Errorf("valkyrie.toml already exists — project already initialized")
		}

		project, _ := cmd.Flags().GetString("project")
		if project == "" {
			project = filepath.Base(cwd)
		}

		serverURL, _ := cmd.Flags().GetString("server")

		// valkyrie.toml
		tomlContent := fmt.Sprintf(`project = %q

[server]
url = %q
token_env = "VALKYRIE_TOKEN"
`, project, serverURL)

		if err := os.WriteFile(tomlPath, []byte(tomlContent), 0644); err != nil {
			return fmt.Errorf("writing valkyrie.toml: %w", err)
		}

		// .mcp.json — Claude Code auto-discovers this
		mcpConfig := `{
  "mcpServers": {
    "valkyrie": {
      "command": "valkyrie",
      "args": ["mcp"]
    }
  }
}
`
		mcpPath := filepath.Join(cwd, ".mcp.json")
		if err := os.WriteFile(mcpPath, []byte(mcpConfig), 0644); err != nil {
			return fmt.Errorf("writing .mcp.json: %w", err)
		}

		// CLAUDE.md — agent instructions
		claudePath := filepath.Join(cwd, "CLAUDE.md")
		if err := os.WriteFile(claudePath, []byte(agentInstructions), 0644); err != nil {
			return fmt.Errorf("writing CLAUDE.md: %w", err)
		}

		// .claude/settings.json — auto-allow valkyrie MCP tools
		claudeDir := filepath.Join(cwd, ".claude")
		if err := os.MkdirAll(claudeDir, 0755); err != nil {
			return fmt.Errorf("creating .claude dir: %w", err)
		}
		settings := `{
  "permissions": {
    "allow": [
      "mcp__valkyrie__*"
    ]
  }
}
`
		settingsPath := filepath.Join(claudeDir, "settings.json")
		if err := os.WriteFile(settingsPath, []byte(settings), 0644); err != nil {
			return fmt.Errorf("writing .claude/settings.json: %w", err)
		}

		fmt.Printf("initialized valkyrie project %q\n", project)
		return nil
	},
}

func init() {
	initCmd.Flags().String("project", "", "Project name (defaults to directory name)")
	initCmd.Flags().String("server", "http://localhost:7777", "Rune server URL")
}

const agentInstructions = `# Valkyrie Rune Agent

You are working in a Valkyrie project. Valkyrie manages function specifications called "runes" on a remote rune server.

## Workflow

1. The user describes what they want to build
2. You refine the requirements with the user until they're ready
3. Use ` + "`requirements_submit`" + ` to push requirements to the rune server
4. The server decomposes requirements into runes, classifies them, and designs specs
5. Use ` + "`requirements_status`" + ` to check progress
6. Present the results to the user for approval
7. Use ` + "`runes_approve`" + ` to commit approved rune specs
8. Use ` + "`runes_reject`" + ` to send feedback on specs that need changes
9. Iterate until the user is satisfied

## Rune Classification

Runes use dot-notation like a standard library:
- Generic (reusable): ` + "`net.http.parse_url`" + `, ` + "`crypto.hash.sha256`" + `, ` + "`text.validate.email`" + `
- Project-specific: ` + "`projectname.auth.validate_token`" + `, ` + "`projectname.payment.calculate_total`" + `

## Available MCP Tools

- ` + "`requirements_submit`" + ` — Push refined requirements to the server
- ` + "`requirements_status`" + ` — Check status of a requirements job
- ` + "`runes_list`" + ` — List runes (by project, namespace, or all)
- ` + "`runes_get`" + ` — Get a rune by fully-qualified name
- ` + "`runes_search`" + ` — Search across the registry
- ` + "`runes_approve`" + ` — Approve/commit a rune spec
- ` + "`runes_reject`" + ` — Reject with feedback (triggers re-design)

## Rune Spec Format

Each rune is a pure function specification:
- **FQN**: Dot-notation name (e.g. ` + "`text.validate.email`" + `)
- **Description**: 1-2 sentences
- **Signature**: Typed function signature using precise types
- **Behavior**: Inputs, outputs, edge cases, constraints
- **Tests**: Positive and negative test cases
`
