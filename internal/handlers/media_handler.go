package handlers

import (
	"context"
	"net/http"
	"sample/internal/database"
	"sample/internal/helpers"
	"sample/internal/utils"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	Dao *database.Queries
	Ctx context.Context
}

func (h *MediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	savedPath := helpers.GenUniquePathAndSave(file, c)

	id, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 0,
			"msg":  "Unauthorized",
		})
		return
	}

	_, uploadErr := h.Dao.UploadMedia(h.Ctx, database.UploadMediaParams{
		UserID:      *id,
		Name:        file.Filename,
		Size:        int32(file.Size),
		ContentType: file.Header.Get("Content-Type"),
		Path:        savedPath,
	})
	if uploadErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 0,
			"msg":  uploadErr.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success",
		"data": gin.H{
			"path": savedPath,
		},
	})
}

func (h *MediaHandler) Uploads(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}

	var paths []string

	id, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 0,
			"msg":  "Unauthorized",
		})
		return
	}

	for _, file := range form.File["files"] {
		savedPath := helpers.GenUniquePathAndSave(file, c)

		_, uploadErr := h.Dao.UploadMedia(h.Ctx, database.UploadMediaParams{
			UserID:      *id,
			Name:        file.Filename,
			Size:        int32(file.Size),
			ContentType: file.Header.Get("Content-Type"),
			Path:        savedPath,
		})
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  uploadErr.Error(),
			})
			return
		}

		paths = append(paths, savedPath)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "Success",
		"data": paths,
	})
}
