package handlers

import (
	"sample/internal/database"
	"sample/internal/helpers"
	"sample/internal/utils"

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

func (h *UserHandler) GetMe(c *gin.Context) {
	id, ok := utils.GetUserID(c)
	if !ok {
		h.handleError(c, "Unauthorized!")
		return
	}

	user, err := h.Dao.GetUser(h.Ctx, *id)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, user, nil)
}

func (h *UserHandler) GetUser(c *gin.Context) {}

func (h *UserHandler) UpdateUser(c *gin.Context) {}

func (h *UserHandler) DeleteUser(c *gin.Context) {}
