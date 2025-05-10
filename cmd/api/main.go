package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"logging/config"
	"logging/internal/common/logger"
	"logging/internal/infra"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load config
	cfg, err := config.LoadAllConfigs("config")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Init logger
	logg := logger.New(cfg.Log)
	logg.Info("⇢ initializing backend...")

	infra, err := infra.InitInfra(ctx, cfg)
	if err != nil {
		logg.Fatalf("infra init error: %v", err)
	}

	engine := SetupRouter(cfg, logg, infra)
	srv := &http.Server{
		Addr:    cfg.App.HostPort,
		Handler: engine,
	}

	// Graceful shutdown
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logg.Infof("⇢ starting HTTP server on %s", cfg.App.HostPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatalf("server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sig
	logg.Info("⇢ shutdown signal received...")

	cancel()
	infra.Close()

	// Timeout to shutdown
	shutdownCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTimeout()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logg.Errorf("graceful shutdown failed: %v", err)
	} else {
		logg.Info("⇢ HTTP server shut down gracefully")
	}
}
