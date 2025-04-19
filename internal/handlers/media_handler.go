package handlers

import (
	"sample/internal/database"
	"sample/internal/helpers"
	"sample/internal/utils"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	BaseHandler
}

func (h *MediaHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	savedPath := helpers.GenUniquePathAndSave(file, c)

	id, ok := utils.GetUserID(c)
	if !ok {
		h.handleError(c, "Unauthorized")
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
		h.handleError(c, uploadErr.Error())
		return
	}

	h.handleSuccess(c, gin.H{
		"path": savedPath,
	}, nil)
}

func (h *MediaHandler) Uploads(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		h.handleError(c, err.Error())
		return
	}

	var paths []string

	id, ok := utils.GetUserID(c)
	if !ok {
		h.handleError(c, "Unauthorized")
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
			h.handleError(c, uploadErr.Error())
			return
		}

		paths = append(paths, savedPath)
	}

	h.handleSuccess(c, paths, nil)
}
