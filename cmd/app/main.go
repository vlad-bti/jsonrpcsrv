package main

import (
	"log"

	"github.com/vlad-bti/jsonrpcsrv/config"
	"github.com/vlad-bti/jsonrpcsrv/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
