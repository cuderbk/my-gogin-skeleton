package main

import (
	"context"
	"log"

	"logging/config"
	"logging/internal/common/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load config
	cfg, err := config.LoadAllConfigs("../../config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	logg := logger.New(cfg.Log)
	logg.Info("⇢ initializing backend...")

	infra, err := InitInfra(ctx, cfg, kafka.TopicHandlerMap{
		"dashboard-logging": kafka.DashboardLoggingHandler(),
		"alert-logging":     kafka.AlertLoggingHandler(),
	})
	if err != nil {
		logg.Fatalf("infra init error: %v", err)
	}
	defer infra.Stop()

	infra.Start(ctx)

	engine := setupRouter(cfg, logg, infra)

	logg.Infof("⇢ starting HTTP server on %s", cfg.App.HostPort)
	if err := engine.Run(cfg.App.HostPort); err != nil {
		logg.Fatalf("server error: %v", err)
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("Shutting down...")
	case err := <-infra.errCh:
		log.Fatalf("Kafka error: %v", err)
	}
}
