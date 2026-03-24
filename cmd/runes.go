package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chrishayen/valkyrie/internal/analyzer"
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
		signature, _ := cmd.Flags().GetString("signature")

		r := runepkg.Rune{
			Name:        name,
			Description: description,
			Signature:   signature,
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
		if cmd.Flags().Changed("signature") {
			r.Signature, _ = cmd.Flags().GetString("signature")
		}
		if cmd.Flags().Changed("version") {
			r.Version, _ = cmd.Flags().GetString("version")
		}

		if !cmd.Flags().Changed("description") && !cmd.Flags().Changed("signature") && !cmd.Flags().Changed("version") {
			return fmt.Errorf("at least one of --description, --signature, or --version is required")
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
		run, err := runner.New(cfg.Agent)
		if err != nil {
			return err
		}

		result, err := hyd.Hydrate(cmd.Context(), args[0], run, os.Stderr)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var runesAnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Decompose requirements into runes via sandbox agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		requirements, _ := cmd.Flags().GetString("requirements")
		yes, _ := cmd.Flags().GetBool("yes")

		run, err := runner.New(cfg.Agent)
		if err != nil {
			return err
		}

		result, err := ana.Analyze(cmd.Context(), requirements, run, os.Stderr)
		if err != nil {
			return err
		}

		printAnalysis(result)

		if len(result.NewRunes) == 0 {
			fmt.Println("\nNo new runes to create.")
			return nil
		}

		if !yes {
			fmt.Print("\nCreate these runes? [y/N] ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		for _, p := range result.NewRunes {
			r := p.ToRune()
			if err := store.Create(r); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create rune %q: %v\n", p.Name, err)
				continue
			}
			fmt.Printf("created rune %q\n", p.Name)
		}

		return nil
	},
}

func printAnalysis(r *analyzer.Result) {
	if len(r.NewRunes) > 0 {
		fmt.Println("=== New runes ===")
		for _, p := range r.NewRunes {
			fmt.Printf("\n  %s\n", p.Name)
			fmt.Printf("    %s\n", p.Description)
			fmt.Printf("    Signature: %s\n", p.Signature)
			fmt.Printf("    Behavior: %s\n", p.Behavior)
			for _, t := range p.PositiveTests {
				fmt.Printf("    + %s\n", t)
			}
			for _, t := range p.NegativeTests {
				fmt.Printf("    - %s\n", t)
			}
		}
		fmt.Println()
	}

	if len(r.ExistingRunes) > 0 {
		fmt.Println("=== Existing runes ===")
		for _, e := range r.ExistingRunes {
			fmt.Printf("\n  %s\n", e.Name)
			fmt.Printf("    Covers: %s\n", e.Covers)
		}
		fmt.Println()
	}
}

func init() {
	runesCmd.AddCommand(runesListCmd)
	runesCmd.AddCommand(runesCreateCmd)
	runesCmd.AddCommand(runesGetCmd)
	runesCmd.AddCommand(runesUpdateCmd)
	runesCmd.AddCommand(runesDeleteCmd)
	runesCmd.AddCommand(runesHydrateCmd)
	runesCmd.AddCommand(runesAnalyzeCmd)

	runesCreateCmd.Flags().String("name", "", "Rune name (slug)")
	runesCreateCmd.Flags().String("description", "", "Rune description")
	runesCreateCmd.Flags().String("signature", "", "Function signature, e.g. (email: string) -> bool")
	_ = runesCreateCmd.MarkFlagRequired("name")
	_ = runesCreateCmd.MarkFlagRequired("description")
	_ = runesCreateCmd.MarkFlagRequired("signature")

	runesUpdateCmd.Flags().String("description", "", "New description")
	runesUpdateCmd.Flags().String("signature", "", "New function signature")
	runesUpdateCmd.Flags().String("version", "", "New version")

	runesAnalyzeCmd.Flags().String("requirements", "", "Plain-text requirements to decompose")
	runesAnalyzeCmd.Flags().Bool("yes", false, "Auto-approve rune creation without prompting")
	_ = runesAnalyzeCmd.MarkFlagRequired("requirements")
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
