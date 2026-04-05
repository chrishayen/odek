package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chrishayen/odek/config"
	"github.com/chrishayen/odek/internal/codegen"
	"github.com/chrishayen/odek/internal/decomposer"
	runepkg "github.com/chrishayen/odek/internal/rune"
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
		if runepkg.IsLeaf(name, storeRuneNames()) {
			codegen.ScaffoldFiles(store.CodeDir(name), runepkg.ShortName(name), config.LangExtension(cfg.Language))
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
			v, _ := cmd.Flags().GetString("version")
			r.Version = runepkg.ParseSemver(v)
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
	Short: "Hydrate a rune (generate code)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := hyd.Hydrate(cmd.Context(), args[0])
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var runesHydrateAllCmd = &cobra.Command{
	Use:   "hydrate-all",
	Short: "Hydrate all un-hydrated runes in parallel",
	RunE: func(cmd *cobra.Command, args []string) error {
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		if concurrency == 0 {
			concurrency = cfg.Concurrency
		}
		verify, _ := cmd.Flags().GetBool("verify")

		result, err := hyd.HydrateAll(cmd.Context(), concurrency, verify, os.Stderr)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

var runesDecomposeCmd = &cobra.Command{
	Use:   "decompose [requirement...]",
	Short: "Decompose requirements into runes",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requirements := strings.Join(args, " ")
		yes, _ := cmd.Flags().GetBool("yes")

		result, err := dec.Decompose(cmd.Context(), requirements, "", os.Stderr, "", "")
		if err != nil {
			return err
		}

		printDecomposition(result)

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

		ext := config.LangExtension(cfg.Language)
		allNames := storeRuneNames()
		for _, p := range result.NewRunes {
			allNames = append(allNames, p.Name)
		}
		for _, p := range result.NewRunes {
			r := p.ToRune()
			if err := store.Create(r); err != nil {
				fmt.Fprintf(os.Stderr, "failed to create rune %q: %v\n", p.Name, err)
				continue
			}
			if runepkg.IsLeaf(p.Name, allNames) {
				codegen.ScaffoldFiles(store.CodeDir(p.Name), runepkg.ShortName(p.Name), ext)
			}
			fmt.Printf("created rune %q\n", p.Name)
		}

		return nil
	},
}

var runesCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for stale references in rune dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		stale, ok, err := store.CheckStaleRefs()
		if err != nil {
			return err
		}
		if stale == 0 {
			fmt.Printf("All %d references up to date.\n", ok)
		} else {
			fmt.Printf("%d stale, %d ok\n", stale, ok)
		}
		return nil
	},
}

var runesVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify hydrated runes against their specs",
	RunE: func(cmd *cobra.Command, args []string) error {
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		if concurrency == 0 {
			concurrency = cfg.Concurrency
		}
		result, err := hyd.VerifyAll(cmd.Context(), concurrency, nil)
		if err != nil {
			return err
		}
		return printJSON(result)
	},
}

func printDecomposition(r *decomposer.Result) {
	if len(r.NewRunes) > 0 {
		fmt.Println("=== New runes ===")
		for _, p := range r.NewRunes {
			fmt.Printf("\n  %s\n", p.Name)
			if p.Description != "" {
				fmt.Printf("    %s\n", p.Description)
			}
			if p.Signature != "" {
				fmt.Printf("    Signature: %s\n", p.Signature)
			}
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
	runesCmd.AddCommand(runesHydrateAllCmd)
	runesCmd.AddCommand(runesDecomposeCmd)
	runesCmd.AddCommand(runesCheckCmd)
	runesCmd.AddCommand(runesVerifyCmd)

	runesCreateCmd.Flags().String("name", "", "Rune name (dot path, e.g. auth.validate_email)")
	runesCreateCmd.Flags().String("description", "", "Rune description")
	runesCreateCmd.Flags().String("signature", "", "Function signature, e.g. (email: string) -> bool")
	_ = runesCreateCmd.MarkFlagRequired("name")
	_ = runesCreateCmd.MarkFlagRequired("description")
	_ = runesCreateCmd.MarkFlagRequired("signature")

	runesUpdateCmd.Flags().String("description", "", "New description")
	runesUpdateCmd.Flags().String("signature", "", "New function signature")
	runesUpdateCmd.Flags().String("version", "", "New version")

	runesDecomposeCmd.Flags().Bool("yes", false, "Auto-approve rune creation without prompting")

	runesHydrateAllCmd.Flags().Int("concurrency", 0, "Max concurrent hydration tasks (default from config)")
	runesHydrateAllCmd.Flags().Bool("verify", false, "Verify implementations after hydration")

	runesVerifyCmd.Flags().Int("concurrency", 0, "Max concurrent verification tasks (default from config)")
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
