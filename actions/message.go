package actions

import (
	"net/http"

	"github.com/FlowerWrong/anychat/services"
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
	uuid := c.Param("uuid")
	room, err := services.FindRoomByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room_uuid": room.Uuid,
	})
}
