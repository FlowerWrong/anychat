package actions

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UploadHandler ...
// Set a lower memory limit for multipart forms (default is 32 MiB)
// router.MaxMultipartMemory = 8 << 20  // 8 MiB
// curl -X POST -H "Content-Type: multipart/form-data" "http://localhost:8080/api/v1/upload" -F "file=@/Users/kingyang/Pictures/gaoyuanyuan.jpg"
func UploadHandler(c *gin.Context) {
	// single file
	file, _ := c.FormFile("file")
	log.Println(file.Filename)

	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, "./tmp/uploader/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("'%s' upload failed! %s", file.Filename, err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("'%s' uploaded!", file.Filename),
	})
}
