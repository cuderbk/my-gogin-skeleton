package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (int, int) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	return page, limit
}
