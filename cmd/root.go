package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/config"
	"github.com/chrishayen/valkyrie/internal/client"
	"github.com/spf13/cobra"
)

var (
	cfg       *config.Config
	apiClient *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "valkyrie",
	Short: "Valkyrie — agentic code orchestration",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for commands that don't need it
		switch cmd.Name() {
		case "init", "serve":
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		apiClient = client.New(cfg.Server.URL, cfg.ResolveToken())
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(mcpCmd)
}
