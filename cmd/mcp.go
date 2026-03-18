package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the Valkyrie MCP server (proxies to HTTP API)",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiURL, _ := cmd.Flags().GetString("api")
		fmt.Printf("Valkyrie MCP server — connecting to API at %s\n", apiURL)
		fmt.Println("MCP server not yet implemented.")
		return nil
	},
}

func init() {
	mcpCmd.Flags().String("api", "http://localhost:8080", "Valkyrie API URL to proxy to")
}
