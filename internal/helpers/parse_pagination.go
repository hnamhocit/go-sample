package helpers

import "github.com/gin-gonic/gin"

type Pagination struct {
	Page  int32 `json:"page" binding:"gte=1"`
	Limit int32 `json:"limit" binding:"gte=1,lte=100"`
}

func ParsePagination(c *gin.Context) (Pagination, error) {
	var pagination Pagination

	if err := c.ShouldBindQuery(&pagination); err != nil {
		return Pagination{}, err
	}

	if pagination.Page < 1 {
		pagination.Page = 1
	}

	if pagination.Limit <= 0 || pagination.Limit > 100 {
		pagination.Limit = 20
	}

	return pagination, nil
}
