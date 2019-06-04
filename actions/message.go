package actions

import (
	"net/http"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/gin-gonic/gin"
)

// SingleChatMsgHandler ...
func SingleChatMsgHandler(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if exists == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current user not found"})
		return
	}
	currentUser := currentUserI.(*models.User)

	lastMsgTimestamp := c.Param("last_msg_timestamp") // FIXME
	withUUID := c.Param("with")
	with, err := services.FindUserByUUID(withUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var ms []models.ChatMessage
	err = db.Engine().Where("((from = ? and to = ?) or (from = ? and to = ?)) and created_at > ?", currentUser.Id, with.Id, with.Id, currentUser.Id, lastMsgTimestamp).Find(&ms)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ms)
}

// RoomChatMsgHandler ...
func RoomChatMsgHandler(c *gin.Context) {
	// currentUserI, exists := c.Get("currentUser")
	// if exists == false {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "current user not found"})
	// 	return
	// }
	// currentUser := currentUserI.(*models.User)

	uuid := c.Param("uuid")
	paramMap := map[string]interface{}{"room_uuid": uuid}
	var roomMessages []models.RoomMessage
	err := db.Engine().SqlTemplateClient("my_room_message_list.tpl", &paramMap).Find(&roomMessages)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, roomMessages)
}
