package cmd

import (
	"context"

	"github.com/chrishayen/odek/internal/feature"
	"github.com/spf13/cobra"
)

var featuresCmd = &cobra.Command{
	Use:   "features",
	Short: "Manage features (top-level rune packages)",
}

var featuresListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all features",
	RunE: func(cmd *cobra.Command, args []string) error {
		features, err := featureStore.List()
		if err != nil {
			return err
		}
		if features == nil {
			features = []feature.Feature{}
		}
		return printJSON(features)
	},
}

var featuresGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a feature by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := featureStore.Get(args[0])
		if err != nil {
			return err
		}
		return printJSON(f)
	},
}

var featuresComposeCmd = &cobra.Command{
	Use:   "compose [name]",
	Short: "Generate dispatcher and wiring code for a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := comp.Compose(context.Background(), args[0])
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

func init() {
	featuresCmd.AddCommand(featuresListCmd)
	featuresCmd.AddCommand(featuresGetCmd)
	featuresCmd.AddCommand(featuresComposeCmd)
}
