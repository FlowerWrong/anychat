package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler ...
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
