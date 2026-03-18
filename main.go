package main

import (
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fmt.Printf("Valkyrie loaded — %d agent(s) configured\n", len(cfg.Agents))
	for name, agent := range cfg.Agents {
		fmt.Printf("  %s: type=%s", name, agent.Type)
		if agent.Model != "" {
			fmt.Printf(" model=%s", agent.Model)
		}
		fmt.Println()
	}
}
