package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/chrishayen/odek/config"
	"github.com/chrishayen/odek/internal/llm"
	"github.com/chrishayen/odek/internal/composer"
	"github.com/chrishayen/odek/internal/decomposer"
	"github.com/chrishayen/odek/internal/feature"
	"github.com/chrishayen/odek/internal/hydrator"
	runepkg "github.com/chrishayen/odek/internal/rune"
	"github.com/chrishayen/odek/internal/validator"
	"github.com/spf13/cobra"
)

var (
	cfg          *config.Config
	store        *runepkg.Store
	featureStore *feature.Store
	client       *llm.Client
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

		providerName, _ := cmd.Flags().GetString("provider")
		prov, err := cfg.ResolveProvider(providerName)
		if err != nil {
			return err
		}

		store = runepkg.NewStore(filepath.Join(cfg.RegistryPath, "runes"), cfg.OutputPath)
		featureStore = feature.NewStore(store, cfg.OutputPath)
		client = llm.New(prov.Model, cfg.Agent.ResolveToken(), cfg.Agent.Mock, prov.Format, prov.BaseURL, prov.MaxTokens)
		val := validator.New(client, cfg.MaxRetries)
		hyd = hydrator.New(store, client, cfg.Language, val)
		dec = decomposer.New(store, client, val)
		comp = composer.New(store, client, cfg.Language)
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
	rootCmd.PersistentFlags().String("provider", "", "Named provider from [providers] in odek.toml")

	rootCmd.AddCommand(runesCmd)
	rootCmd.AddCommand(featuresCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(serveCmd)
}
