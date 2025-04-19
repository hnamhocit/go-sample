package handlers

import (
	"context"
	"net/http"
	"sample/internal/database"
	"sample/internal/dtos"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	Dao *database.Queries
	Ctx context.Context
}

func (h *BaseHandler) handleSuccess(c *gin.Context, data interface{}, msg *string) {
	if msg != nil {
		c.JSON(http.StatusOK, dtos.Response{
			Code: 1,
			Msg:  *msg,
			Data: data,
		})
		return
	}

	c.JSON(http.StatusOK, dtos.Response{
		Code: 1,
		Msg:  "Success",
		Data: data,
	})
	return
}

func (h *BaseHandler) handleError(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, dtos.Response{
		Code: 0,
		Msg:  msg,
		Data: nil,
	})
	return
}
