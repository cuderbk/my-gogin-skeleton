package middleware

import (
	"logging/internal/common/response"

	"github.com/gin-gonic/gin"
)

func ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() != "application/json" {
			response.BadRequest(c, "Content-Type must be application/json", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
