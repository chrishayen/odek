package cmd

import (
	"context"
	"fmt"

	"github.com/chrishayen/valkyrie/internal/feature"
	"github.com/spf13/cobra"
)

var featuresCmd = &cobra.Command{
	Use:   "features",
	Short: "Manage features in the registry",
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

var featuresCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new feature",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if err := featureStore.Create(name, description); err != nil {
			return err
		}
		created, err := featureStore.Get(name)
		if err != nil {
			return err
		}
		return printJSON(created)
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

var featuresUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := featureStore.Get(args[0])
		if err != nil {
			return err
		}

		changed := false
		if cmd.Flags().Changed("version") {
			f.Version, _ = cmd.Flags().GetString("version")
			changed = true
		}
		if cmd.Flags().Changed("status") {
			f.Status, _ = cmd.Flags().GetString("status")
			changed = true
		}

		if !changed {
			return fmt.Errorf("at least one of --version or --status is required")
		}

		if err := featureStore.Update(*f); err != nil {
			return err
		}
		updated, err := featureStore.Get(args[0])
		if err != nil {
			return err
		}
		return printJSON(updated)
	},
}

var featuresDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a feature",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := featureStore.Delete(args[0]); err != nil {
			return err
		}
		fmt.Printf("feature %q deleted\n", args[0])
		return nil
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
	featuresCmd.AddCommand(featuresCreateCmd)
	featuresCmd.AddCommand(featuresGetCmd)
	featuresCmd.AddCommand(featuresUpdateCmd)
	featuresCmd.AddCommand(featuresDeleteCmd)
	featuresCmd.AddCommand(featuresComposeCmd)

	featuresCreateCmd.Flags().String("name", "", "Feature name (slug)")
	featuresCreateCmd.Flags().String("description", "", "Feature description")
	_ = featuresCreateCmd.MarkFlagRequired("name")
	_ = featuresCreateCmd.MarkFlagRequired("description")

	featuresUpdateCmd.Flags().String("version", "", "New version")
	featuresUpdateCmd.Flags().String("status", "", "New status (draft, reviewed, stable)")
}
