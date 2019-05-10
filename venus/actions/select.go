package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SelectHandler 返回websocket连接的主机和port
func SelectHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"host": "1.1.1.1",
	})
}
