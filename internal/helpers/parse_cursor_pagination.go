package helpers

import "github.com/gin-gonic/gin"

type CursorPagination struct {
	Cursor int32 `json:"cursor" binding:"gte=0"`
	Limit  int32 `json:"limit" binding:"gte=1,lte=100"`
}

func ParseCursorPagination(c *gin.Context) (CursorPagination, error) {
	var pagination CursorPagination

	if err := c.ShouldBindQuery(&pagination); err != nil {
		return CursorPagination{}, err
	}

	if pagination.Cursor < 0 {
		pagination.Cursor = 0
	}
	if pagination.Limit <= 0 || pagination.Limit > 100 {
		pagination.Limit = 20
	}

	return pagination, nil
}
