//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"logging/config"
	"logging/internal/common/logger"
	"logging/internal/common/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProvideLogConfig(cfg *config.Config) config.LogConfig {
	return cfg.Log
}

func initializeServer(configPath string) (*gin.Engine, error) {
	wire.Build(
		config.LoadAllConfigs,
		ProvideLogConfig,

		logger.New,

		buildRouter,
	)
	return &gin.Engine{}, nil
}

func buildRouter(
	cfg *config.Config,
	log *zap.SugaredLogger,

) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(), middleware.RequestID())

	api := r.Group("/api/v0")

	api.Use(middleware.ValidateContentType())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Infof("starting HTTP server on %s", cfg.App.HostPort)
	if err := r.Run(cfg.App.HostPort); err != nil {
		log.Fatalf("server error: %v", err)
	}

	return r
}
