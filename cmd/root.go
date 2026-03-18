package cmd

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/internal/config"
	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "valkyrie",
	Short: "The rune server — agentic code orchestration",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	var err error
	cfg, err = config.Load("valkyrie.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}
}
