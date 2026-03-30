package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/chrishayen/valkyrie/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.New(), tea.WithAltScreen())
		_, err := p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
