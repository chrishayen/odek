package cmd

import (
	"fmt"

	"github.com/chrishayen/odek/internal/app"
	"github.com/spf13/cobra"
)

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Manage apps in the registry",
}

var appsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		apps, err := appStore.List()
		if err != nil {
			return err
		}
		if apps == nil {
			apps = []app.App{}
		}
		return printJSON(apps)
	},
}

var appsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new app",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if err := appStore.Create(name, description); err != nil {
			return err
		}
		created, err := appStore.Get(name)
		if err != nil {
			return err
		}
		return printJSON(created)
	},
}

var appsGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get an app by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := appStore.Get(args[0])
		if err != nil {
			return err
		}
		return printJSON(a)
	},
}

var appsUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := appStore.Get(args[0])
		if err != nil {
			return err
		}

		changed := false
		if cmd.Flags().Changed("version") {
			a.Version, _ = cmd.Flags().GetString("version")
			changed = true
		}
		if cmd.Flags().Changed("status") {
			a.Status, _ = cmd.Flags().GetString("status")
			changed = true
		}
		if cmd.Flags().Changed("entry-point") {
			a.EntryPoint, _ = cmd.Flags().GetString("entry-point")
			changed = true
		}

		if !changed {
			return fmt.Errorf("at least one of --version, --status, or --entry-point is required")
		}

		if err := appStore.Update(*a); err != nil {
			return err
		}
		updated, err := appStore.Get(args[0])
		if err != nil {
			return err
		}
		return printJSON(updated)
	},
}

var appsDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := appStore.Delete(args[0]); err != nil {
			return err
		}
		fmt.Printf("app %q deleted\n", args[0])
		return nil
	},
}

func init() {
	appsCmd.AddCommand(appsListCmd)
	appsCmd.AddCommand(appsCreateCmd)
	appsCmd.AddCommand(appsGetCmd)
	appsCmd.AddCommand(appsUpdateCmd)
	appsCmd.AddCommand(appsDeleteCmd)

	appsCreateCmd.Flags().String("name", "", "App name (slug)")
	appsCreateCmd.Flags().String("description", "", "App description")
	_ = appsCreateCmd.MarkFlagRequired("name")
	_ = appsCreateCmd.MarkFlagRequired("description")

	appsUpdateCmd.Flags().String("version", "", "New version")
	appsUpdateCmd.Flags().String("status", "", "New status (draft, reviewed, stable)")
	appsUpdateCmd.Flags().String("entry-point", "", "Entry point feature")
}
