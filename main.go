package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chrishayen/valkyrie/config"
	"github.com/chrishayen/valkyrie/internal/runner"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	// Usage: valkyrie <agent-name> <task>
	if len(os.Args) == 3 {
		agentName := os.Args[1]
		task := os.Args[2]

		agent, ok := cfg.Agents[agentName]
		if !ok {
			fmt.Fprintf(os.Stderr, "error: agent %q not found in config\n", agentName)
			os.Exit(1)
		}

		r, err := runner.New(agent)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		out, err := r.Run(context.Background(), task)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(out)
		return
	}

	// Default: print loaded config summary
	fmt.Printf("Valkyrie loaded — %d agent(s) configured\n", len(cfg.Agents))
	for name, agent := range cfg.Agents {
		fmt.Printf("  %s: type=%s", name, agent.Type)
		if agent.Model != "" {
			fmt.Printf(" model=%s", agent.Model)
		}
		fmt.Println()
	}
}
