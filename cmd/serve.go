package cmd

import (
	"fmt"
	"net/http"

	"github.com/chrishayen/valkyrie/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP API server",
	RunE: func(cmd *cobra.Command, args []string) error {
		cancelProxy, err := startProxy(false)
		if err != nil {
			return fmt.Errorf("proxy: %w", err)
		}
		defer cancelProxy()

		port, _ := cmd.Flags().GetInt("port")
		if port == 0 {
			port = cfg.Server.Port
		}

		s := server.New(cfg, store, featureStore, appStore, dec, hyd)
		addr := fmt.Sprintf(":%d", port)
		fmt.Printf("valkyrie api server listening on %s\n", addr)
		return http.ListenAndServe(addr, s)
	},
}

func init() {
	serveCmd.Flags().Int("port", 0, "Port to listen on (default from config)")
}
