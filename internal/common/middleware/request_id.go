package middleware

import (
	"logging/internal/common/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			uuid, err := uuid.NewUUID()
			if err != nil {
				response.WrapError(c, http.StatusInternalServerError, "Failed to generate request ID", err.Error())
				c.Abort()
				return
			}
			reqID = uuid.String()
		}
		c.Set("request_id", reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Next()
	}
}
