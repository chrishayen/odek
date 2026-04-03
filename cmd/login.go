package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/router-for-me/CLIProxyAPI/v6/sdk/auth"
	"github.com/router-for-me/CLIProxyAPI/v6/sdk/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Claude via OAuth",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := os.Getenv("CLIPROXY_CONFIG")
		if configPath == "" {
			configPath = "config.yaml"
		}

		proxyCfg, err := config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("load proxy config: %w", err)
		}

		cancelProxy, err := startProxy(false, true)
		if err != nil {
			return fmt.Errorf("proxy: %w", err)
		}
		defer cancelProxy()

		store := auth.NewFileTokenStore()
		store.SetBaseDir(proxyCfg.AuthDir)
		mgr := auth.NewManager(store, auth.NewClaudeAuthenticator())

		_, _, err = mgr.Login(context.Background(), "claude", proxyCfg, nil)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
