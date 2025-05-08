package main

import (
	"logging/config"
	"logging/internal/common/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupRouter(cfg *config.Config, logg *zap.SugaredLogger, infra *Infra) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(), middleware.RequestID())

	api := r.Group("/api/v0")
	api.Use(middleware.ValidateContentType())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Bạn có thể dùng: infra.DB, infra.Redis

	return r
}
