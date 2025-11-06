package main

import (
	"fmt"
	"os"

	"svrn/internal/config"
	"svrn/internal/logging"
	"svrn/pkg/agent"
)

var version = "0.0.0-dev"

func main() {
	// Early flag check for --version without loading full config
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "version" || arg == "-v" {
			fmt.Println("svrn", version)
			return
		}
	}

	// Initialize logging (temporary stderr until config loads)
	log := logging.New()

	// Load config from flags/env/YAML with proper precedence
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("config load failed", "error", err)
	}

	// Create agent with resolved configuration
	ag, err := agent.New(cfg)
	if err != nil {
		log.Fatal("agent init failed", "error", err)
	}

	// Start node runtime
	if err := ag.Start(); err != nil {
		log.Fatal("agent start failed", "error", err)
	}
}
