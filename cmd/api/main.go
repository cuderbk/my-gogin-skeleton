package main

import (
	"context"
	"log"

	"logging/config"
	"logging/internal/common/logger"
)

func main() {
	ctx := context.Background()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load config
	cfg, err := config.LoadAllConfigs("../../config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	logg := logger.New(cfg.Log)
	logg.Info("⇢ initializing backend...")

	infra, err := InitInfra(ctx, cfg)
	if err != nil {
		logg.Fatalf("infra init error: %v", err)
	}
	defer infra.Close()

	engine := setupRouter(cfg, logg, infra)

	logg.Infof("⇢ starting HTTP server on %s", cfg.App.HostPort)
	if err := engine.Run(cfg.App.HostPort); err != nil {
		logg.Fatalf("server error: %v", err)
	}
}
