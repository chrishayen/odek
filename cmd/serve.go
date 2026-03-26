package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/internal/adaptor"
	"github.com/chrishayen/valkyrie/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the rune server",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, _ := cmd.Flags().GetInt("port")
		dataDir, _ := cmd.Flags().GetString("data-dir")
		token, _ := cmd.Flags().GetString("token")
		tokenEnv, _ := cmd.Flags().GetString("token-env")
		apiKeyEnv, _ := cmd.Flags().GetString("api-key-env")
		model, _ := cmd.Flags().GetString("model")

		if token == "" && tokenEnv != "" {
			token = os.Getenv(tokenEnv)
		}

		apiKey := os.Getenv(apiKeyEnv)
		if apiKey == "" {
			return fmt.Errorf("Anthropic API key not set — set %s environment variable", apiKeyEnv)
		}

		a := adaptor.NewClaude(apiKey, model)
		s := server.New(dataDir, token, a)
		return s.ListenAndServe(server.Addr(port))
	},
}

func init() {
	home, _ := os.UserHomeDir()
	defaultData := home + "/.valkyrie/data"

	serveCmd.Flags().Int("port", 7777, "Port to listen on")
	serveCmd.Flags().String("data-dir", defaultData, "Data directory for rune storage")
	serveCmd.Flags().String("token", "", "Bearer token for authentication")
	serveCmd.Flags().String("token-env", "VALKYRIE_TOKEN", "Environment variable containing the bearer token")
	serveCmd.Flags().String("api-key-env", "ANTHROPIC_API_KEY", "Environment variable containing the Anthropic API key")
	serveCmd.Flags().String("model", "claude-sonnet-4-5-20250514", "Claude model to use for rune processing")
}
