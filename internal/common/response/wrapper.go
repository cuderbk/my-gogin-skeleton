package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StandardResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// WrapSuccess returns a success response with optional data
func WrapSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// WrapCreated returns a 201 Created response
func WrapCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, StandardResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// WrapError returns an error response with custom status code
func WrapError(c *gin.Context, status int, message string, err interface{}) {
	c.JSON(status, StandardResponse{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func Unauthorized(c *gin.Context, msg string) {
	WrapError(c, http.StatusUnauthorized, msg, nil)
}

// Shortcut helpers
func BadRequest(c *gin.Context, msg string, err interface{}) {
	WrapError(c, http.StatusBadRequest, msg, err)
}

func NotFound(c *gin.Context, msg string) {
	WrapError(c, http.StatusNotFound, msg, nil)
}

func InternalServer(c *gin.Context, msg string, err interface{}) {
	WrapError(c, http.StatusInternalServerError, msg, err)
}
