package cmd

import (
	"encoding/json"
	"fmt"

	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/chrishayen/valkyrie/internal/runner"
	"github.com/spf13/cobra"
)

var runesCmd = &cobra.Command{
	Use:   "runes",
	Short: "Manage runes in the registry",
}

var runesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all runes",
	RunE: func(cmd *cobra.Command, args []string) error {
		runes, err := store.List()
		if err != nil {
			return err
		}
		if runes == nil {
			runes = []runepkg.Rune{}
		}
		return printJSON(runes)
	},
}

var runesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new rune",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		r := runepkg.Rune{
			Name:        name,
			Description: description,
		}
		if err := store.Create(r); err != nil {
			return err
		}
		created, err := store.Get(name)
		if err != nil {
			return err
		}
		return printJSON(created)
	},
}

var runesGetCmd = &cobra.Command{
	Use:   "get [name]",
	Short: "Get a rune by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := store.Get(args[0])
		if err != nil {
			return err
		}
		return printJSON(r)
	},
}

var runesUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a rune",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		r, err := store.Get(args[0])
		if err != nil {
			return err
		}

		if cmd.Flags().Changed("description") {
			r.Description, _ = cmd.Flags().GetString("description")
		}
		if cmd.Flags().Changed("version") {
			r.Version, _ = cmd.Flags().GetString("version")
		}

		if !cmd.Flags().Changed("description") && !cmd.Flags().Changed("version") {
			return fmt.Errorf("at least one of --description or --version is required")
		}

		if err := store.Update(*r); err != nil {
			return err
		}
		updated, err := store.Get(args[0])
		if err != nil {
			return err
		}
		return printJSON(updated)
	},
}

var runesDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a rune",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := store.Delete(args[0]); err != nil {
			return err
		}
		fmt.Printf("rune %q deleted\n", args[0])
		return nil
	},
}

var runesHydrateCmd = &cobra.Command{
	Use:   "hydrate [name]",
	Short: "Hydrate a rune (generate code via sandbox agent)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sandbox, _ := cmd.Flags().GetString("sandbox")

		if sandbox == "" {
			if len(cfg.Agents) == 0 {
				return fmt.Errorf("no agents configured — add an agent to your config")
			}
			for name := range cfg.Agents {
				sandbox = name
				break
			}
		}

		agent, ok := cfg.Agents[sandbox]
		if !ok {
			return fmt.Errorf("sandbox %q not found in config", sandbox)
		}

		run, err := runner.New(agent)
		if err != nil {
			return err
		}

		result, err := hyd.Hydrate(cmd.Context(), args[0], run)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

func init() {
	runesCmd.AddCommand(runesListCmd)
	runesCmd.AddCommand(runesCreateCmd)
	runesCmd.AddCommand(runesGetCmd)
	runesCmd.AddCommand(runesUpdateCmd)
	runesCmd.AddCommand(runesDeleteCmd)
	runesCmd.AddCommand(runesHydrateCmd)

	runesCreateCmd.Flags().String("name", "", "Rune name (slug)")
	runesCreateCmd.Flags().String("description", "", "Rune description")
	_ = runesCreateCmd.MarkFlagRequired("name")
	_ = runesCreateCmd.MarkFlagRequired("description")

	runesUpdateCmd.Flags().String("description", "", "New description")
	runesUpdateCmd.Flags().String("version", "", "New version")

	runesHydrateCmd.Flags().String("sandbox", "", "Sandbox agent name from config (defaults to first agent)")
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
