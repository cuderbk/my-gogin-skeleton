package response

import (
	"github.com/gin-gonic/gin"
)

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

func WrapPagination(c *gin.Context, message string, items interface{}, page, limit, total int) {
	WrapSuccess(c, message, gin.H{
		"items": items,
		"pagination": PaginationMeta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	})
}
