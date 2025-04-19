package utils

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (*int32, bool) {
	userId, ok := c.Get("user_id")
	if !ok {
		return nil, false
	}

	id, ok := userId.(int32)
	if !ok {
		return nil, false
	}

	return &id, true
}
