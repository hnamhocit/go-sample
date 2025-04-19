package handlers

import (
	"context"
	"sample/internal/database"
	"sample/internal/dtos"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Dao *database.Queries
	Ctx context.Context
}

func (r *UserHandler) GetUsers(c *gin.Context) {
	var pagination dtos.PaginationDTO

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(400, dtos.Response{
			Code: 0,
			Msg:  "Bad Request",
		})
		return
	}

	users, err := r.Dao.GetUsers(c, database.GetUsersParams{
		Limit:  int32(pagination.Size),
		Offset: int32((pagination.Page - 1) * pagination.Size),
	})

	if err != nil {
		c.JSON(500, dtos.Response{
			Code: 0,
			Msg:  "Internal Server Error",
		})
		return
	}

	c.JSON(200, dtos.Response{
		Code: 1,
		Msg:  "Success",
		Data: users,
	})
}

func (r *UserHandler) GetUser(c *gin.Context) {}

func (r *UserHandler) UpdateUser(c *gin.Context) {}

func (r *UserHandler) DeleteUser(c *gin.Context) {}
