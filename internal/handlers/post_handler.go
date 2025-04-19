package handlers

import (
	"sample/internal/database"
	"sample/internal/helpers"
	"sample/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	BaseHandler
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	pagination, err := helpers.ParseCursorPagination(c)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	posts, err := h.Dao.GetPosts(h.Ctx, database.GetPostsParams{
		ID:    pagination.Cursor,
		Limit: pagination.Limit,
	})
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, gin.H{"posts": posts}, nil)
}

type CreatePostDTO struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var dto CreatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleError(c, err.Error())
		return
	}

	id, ok := utils.GetUserID(c)
	if !ok {
		h.handleError(c, "User not found")
		return
	}

	post, err := h.Dao.CreatePost(h.Ctx, database.CreatePostParams{
		Title:   dto.Title,
		Content: dto.Content,
		UserID:  *id,
	})
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, gin.H{"post": post}, nil)
}

type UpdatePostDTO struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	var dto UpdatePostDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		h.handleError(c, err.Error())
		return
	}

	id := c.Param("id")

	_id, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	post, err := h.Dao.UpdatePost(h.Ctx, database.UpdatePostParams{
		ID:      int32(_id),
		Title:   dto.Title,
		Content: dto.Content,
	})
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, gin.H{"post": post}, nil)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	_id, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	_, err = h.Dao.DeletePost(h.Ctx, int32(_id))
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, true, nil)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id := c.Param("id")

	_id, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	post, err := h.Dao.GetPost(h.Ctx, int32(_id))
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	h.handleSuccess(c, post, nil)
}
