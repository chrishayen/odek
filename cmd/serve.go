package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/chrishayen/valkyrie/internal/api"
	runepkg "github.com/chrishayen/valkyrie/internal/rune"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Valkyrie HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		addr, _ := cmd.Flags().GetString("addr")

		token := cfg.Auth.ResolveToken()
		if !cfg.Auth.Disabled && token == "" {
			return fmt.Errorf("auth.token is not set — set it in config, export VALKYRIE_API_TOKEN, or set auth.disabled = true for local use")
		}

		store := runepkg.NewStore(cfg.RegistryPath)
		server := api.NewServer(store, api.Options{
			AuthToken:    token,
			AuthDisabled: cfg.Auth.Disabled,
		})

		if cfg.Auth.Disabled {
			fmt.Fprintln(os.Stdout, "⚠️  Auth disabled — for local use only")
		}
		fmt.Fprintf(os.Stdout, "Valkyrie serving on %s\n", addr)
		fmt.Fprintf(os.Stdout, "Registry: %s\n", cfg.RegistryPath)
		return http.ListenAndServe(addr, server)
	},
}

func init() {
	serveCmd.Flags().String("addr", ":8080", "Address to listen on")
}
