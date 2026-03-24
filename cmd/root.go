package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/config"
	"github.com/chrishayen/valkyrie/internal/analyzer"
	"github.com/chrishayen/valkyrie/internal/composer"
	"github.com/chrishayen/valkyrie/internal/feature"
	"github.com/chrishayen/valkyrie/internal/hydrator"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/spf13/cobra"
)

var (
	cfg          *config.Config
	store        *runepkg.Store
	featureStore *feature.Store
	hyd          *hydrator.Hydrator
	ana          *analyzer.Analyzer
	comp         *composer.Composer
)

var rootCmd = &cobra.Command{
	Use:   "valkyrie",
	Short: "Valkyrie — agentic code orchestration",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for init — config doesn't exist yet
		if cmd.Name() == "init" {
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		store = runepkg.NewStore(cfg.RegistryPath)
		featureStore = feature.NewStore(cfg.RegistryPath)
		hyd = hydrator.New(store)
		ana = analyzer.New(store)
		comp = composer.New(featureStore, store)
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
	rootCmd.AddCommand(featuresCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(mcpCmd)
}
