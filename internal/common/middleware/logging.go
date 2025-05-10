package middleware

import (
	"time"

	"my-gogin-skeleton/internal/common/logger"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func ZapLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		// Try to get request ID if set by previous middleware
		reqID, _ := c.Get(RequestIDKey)

		logger.Logger.Infow("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency.String(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"request_id", reqID,
		)
	}
}
