package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/config"
	"github.com/chrishayen/valkyrie/internal/hydrator"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/spf13/cobra"
)

var (
	cfg   *config.Config
	store *runepkg.Store
	hyd   *hydrator.Hydrator
)

var rootCmd = &cobra.Command{
	Use:   "valkyrie",
	Short: "Valkyrie — agentic code orchestration",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		store = runepkg.NewStore(cfg.RegistryPath)
		hyd = hydrator.New(store)
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
	rootCmd.AddCommand(runesCmd)
}
