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

		if !cfg.Auth.Disabled {
			fmt.Fprintln(os.Stderr, "warning: auth is not disabled — set auth.disabled = true in config for local use")
		}

		store := runepkg.NewStore(cfg.RegistryPath)
		server := api.NewServer(store, cfg.Agents)

		fmt.Fprintf(os.Stdout, "Valkyrie serving on %s\n", addr)
		fmt.Fprintf(os.Stdout, "Registry: %s\n", cfg.RegistryPath)
		return http.ListenAndServe(addr, server)
	},
}

func init() {
	serveCmd.Flags().String("addr", ":8080", "Address to listen on")
}
