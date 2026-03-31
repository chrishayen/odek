package cmd

import (
	"fmt"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chrishayen/odek/internal/chat"
	"github.com/chrishayen/odek/internal/server"
	"github.com/chrishayen/odek/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		cancelProxy, err := startProxy(true)
		if err != nil {
			return fmt.Errorf("proxy: %w", err)
		}
		defer cancelProxy()

		port := cfg.Server.Port
		s := server.New(cfg, store, featureStore, appStore, dec, hyd)
		go http.ListenAndServe(fmt.Sprintf(":%d", port), s)

		chatStore := chat.NewStore(cfg.RegistryPath)
		p := tea.NewProgram(tui.New(port, cfg.RegistryPath, chatStore), tea.WithAltScreen())
		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
