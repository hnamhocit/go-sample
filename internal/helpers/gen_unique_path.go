package helpers

import (
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GenUniquePathAndSave(file *multipart.FileHeader, c *gin.Context) string {
	parts := strings.Split(file.Filename, ".")
	filePath := fmt.Sprintf("assets/media/%v-%v.%v", parts[0], time.Now().Unix(), parts[len(parts)-1])
	c.SaveUploadedFile(file, filePath)
	return filePath
}
