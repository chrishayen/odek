package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chrishayen/valkyrie/framework"
	"github.com/chrishayen/valkyrie/internal/analyzer"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Valkyrie project in the current directory",
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

		// valkyrie.toml
		tomlContent := fmt.Sprintf(`project = %q
# output_path = "src"

[agent]
type = "claude-sub"
# model = "claude-sonnet-4-5"
# token parses from ~/.claude/.credentials.json by default
`, project)

		if err := os.WriteFile(tomlPath, []byte(tomlContent), 0644); err != nil {
			return fmt.Errorf("writing valkyrie.toml: %w", err)
		}

		// runes/
		runesDir := filepath.Join(cwd, "runes")
		if err := os.MkdirAll(runesDir, 0755); err != nil {
			return fmt.Errorf("creating runes dir: %w", err)
		}

		// src/ — hydrated code output
		srcDir := filepath.Join(cwd, "src")
		if err := os.MkdirAll(srcDir, 0755); err != nil {
			return fmt.Errorf("creating src dir: %w", err)
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

		// CLAUDE.md — rune-agent instructions, auto-loaded by Claude Code
		claudePath := filepath.Join(cwd, "CLAUDE.md")
		if err := os.WriteFile(claudePath, []byte(analyzer.Instructions), 0644); err != nil {
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

		// internal/dispatch/ — prebuilt dispatch framework
		dispatchDir := filepath.Join(cwd, "internal", "dispatch")
		if err := os.MkdirAll(dispatchDir, 0755); err != nil {
			return fmt.Errorf("creating dispatch dir: %w", err)
		}
		dispatchPath := filepath.Join(dispatchDir, "dispatch.go")
		if err := os.WriteFile(dispatchPath, []byte(framework.GoDispatch), 0644); err != nil {
			return fmt.Errorf("writing dispatch.go: %w", err)
		}

		fmt.Printf("initialized valkyrie project %q\n", project)
		return nil
	},
}

func init() {
	initCmd.Flags().String("project", "", "Project name (defaults to directory name)")
}
