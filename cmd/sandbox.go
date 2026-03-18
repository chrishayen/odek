package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/internal/sandbox"
	"github.com/spf13/cobra"
)

var sandboxCmd = &cobra.Command{
	Use:   "sandbox",
	Short: "Manage sandbox configurations",
}

var sandboxCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new sandbox configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		stype, _ := cmd.Flags().GetString("type")
		model, _ := cmd.Flags().GetString("model")
		apiKeyEnv, _ := cmd.Flags().GetString("api-key-env")
		image, _ := cmd.Flags().GetString("image")

		if name == "" || stype == "" {
			return fmt.Errorf("--name and --type are required")
		}

		s := sandbox.Sandbox{
			Name:      name,
			Type:      stype,
			Model:     model,
			APIKeyEnv: apiKeyEnv,
			Image:     image,
		}

		if err := sandbox.Create(cfg.RegistryPath, s); err != nil {
			return err
		}
		fmt.Printf("✓ Sandbox %q created (type: %s)\n", name, stype)
		return nil
	},
}

var sandboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sandbox configurations",
	RunE: func(cmd *cobra.Command, args []string) error {
		sandboxes, err := sandbox.List(cfg.RegistryPath)
		if err != nil {
			return err
		}
		if len(sandboxes) == 0 {
			fmt.Println("No sandboxes configured.")
			return nil
		}
		fmt.Printf("%-20s %-15s %s\n", "NAME", "TYPE", "MODEL/IMAGE")
		fmt.Printf("%-20s %-15s %s\n", "----", "----", "-----------")
		for _, s := range sandboxes {
			detail := s.Model
			if detail == "" {
				detail = s.Image
			}
			fmt.Printf("%-20s %-15s %s\n", s.Name, s.Type, detail)
		}
		return nil
	},
}

var sandboxGetCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Get a sandbox configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := sandbox.Get(cfg.RegistryPath, args[0])
		if err != nil {
			return err
		}
		fmt.Printf("name      = %q\n", s.Name)
		fmt.Printf("type      = %q\n", s.Type)
		if s.Model != "" {
			fmt.Printf("model     = %q\n", s.Model)
		}
		if s.APIKeyEnv != "" {
			fmt.Printf("api_key_env = %q\n", s.APIKeyEnv)
		}
		if s.Image != "" {
			fmt.Printf("image     = %q\n", s.Image)
		}
		return nil
	},
}

var sandboxDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a sandbox configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sandbox.Delete(cfg.RegistryPath, args[0]); err != nil {
			return err
		}
		fmt.Printf("✓ Sandbox %q deleted\n", args[0])
		return nil
	},
}

func init() {
	sandboxCreateCmd.Flags().String("name", "", "Sandbox name (required)")
	sandboxCreateCmd.Flags().String("type", "", "Sandbox type: claude-api, claude-max, docker (required)")
	sandboxCreateCmd.Flags().String("model", "", "Model to use (claude-api, claude-max)")
	sandboxCreateCmd.Flags().String("api-key-env", "", "Env var containing the API key (claude-api)")
	sandboxCreateCmd.Flags().String("image", "", "Docker image (docker)")

	sandboxCmd.AddCommand(sandboxCreateCmd)
	sandboxCmd.AddCommand(sandboxListCmd)
	sandboxCmd.AddCommand(sandboxGetCmd)
	sandboxCmd.AddCommand(sandboxDeleteCmd)
	rootCmd.AddCommand(sandboxCmd)

	_ = os.MkdirAll("registry/sandboxes", 0755)
}
