package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/odek/config"
	"github.com/chrishayen/odek/internal/app"
	"github.com/chrishayen/odek/internal/claude"
	"github.com/chrishayen/odek/internal/composer"
	"github.com/chrishayen/odek/internal/decomposer"
	"github.com/chrishayen/odek/internal/feature"
	"github.com/chrishayen/odek/internal/hydrator"
	runepkg "github.com/chrishayen/odek/internal/rune"
	"github.com/spf13/cobra"
)

var (
	cfg          *config.Config
	store        *runepkg.Store
	featureStore *feature.Store
	appStore     *app.Store
	client       *claude.Client
	hyd          *hydrator.Hydrator
	dec          *decomposer.Decomposer
	comp         *composer.Composer
)

var rootCmd = &cobra.Command{
	Use:   "odek",
	Short: "Odek — Tree Composition CLI and Rune Server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "init" {
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("config: %w", err)
		}
		store = runepkg.NewStore(cfg.RegistryPath, cfg.OutputPath)
		featureStore = feature.NewStore(cfg.RegistryPath, cfg.OutputPath)
		appStore = app.NewStore(cfg.RegistryPath, cfg.OutputPath)
		client = claude.New(cfg.Agent.Model, cfg.Agent.ResolveToken(), cfg.Agent.Mock)
		hyd = hydrator.New(store, client, cfg.Language)
		dec = decomposer.New(store, client, cfg.Project)
		comp = composer.New(featureStore, store, client, cfg.Language)
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
	rootCmd.AddCommand(appsCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(serveCmd)
}
