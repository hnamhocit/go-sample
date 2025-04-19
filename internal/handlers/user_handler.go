package handlers

import (
	"sample/internal/database"
	"sample/internal/helpers"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	BaseHandler
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	pagination, err := helpers.ParsePagination(c)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	users, err := h.Dao.GetUsers(c, database.GetUsersParams{
		Limit:  int32(pagination.Limit),
		Offset: int32((pagination.Page - 1) * pagination.Limit),
	})

	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, users, nil)
}

func (h *UserHandler) GetUser(c *gin.Context) {}

func (h *UserHandler) UpdateUser(c *gin.Context) {}

func (h *UserHandler) DeleteUser(c *gin.Context) {}
