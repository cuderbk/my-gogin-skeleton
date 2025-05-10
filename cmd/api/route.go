package main

import (
	"my-gogin-skeleton/config"
	"my-gogin-skeleton/internal/common/middleware"
	"my-gogin-skeleton/internal/infra"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(cfg *config.Config, logg *zap.SugaredLogger, infra *infra.Infra) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(), middleware.RequestID())

	api := r.Group("/api/v0")
	api.Use(middleware.ValidateContentType())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// infra.DB, infra.Redis

	return r
}
