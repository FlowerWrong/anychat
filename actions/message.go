package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SingleChatMsgHandler ...
func SingleChatMsgHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// RoomChatMsgHandler ...
func RoomChatMsgHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
