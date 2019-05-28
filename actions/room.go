package actions

import (
	"log"
	"net/http"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/FlowerWrong/util"
	"github.com/gin-gonic/gin"
)

type roomVM struct {
	UserUUIDs []string `form:"user_uuids" json:"user_uuids" xml:"user_uuids" binding:"required"`
	Name      string   `form:"name" json:"name" xml:"name" binding:"required"`
	Intro     string   `form:"intro" json:"intro" xml:"intro" binding:"required"`
	Logo      string   `form:"logo" json:"logo" xml:"logo" binding:"required"`
}

// CreateRoomHandler ...
func CreateRoomHandler(c *gin.Context) {
	var roomForm roomVM
	if err := c.ShouldBindJSON(&roomForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(roomForm)

	room := new(models.Room)
	room.Uuid = util.UUID()
	room.Name = roomForm.Name
	room.Intro = roomForm.Intro
	room.Logo = roomForm.Logo
	err := utils.InsertRecord(room)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, uuid := range roomForm.UserUUIDs {
		roomUser := new(models.RoomUser)
		roomUser.Uuid = util.UUID()
		roomUser.RoomId = room.Id

		user, err := services.FindUserByUUID(uuid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		roomUser.UserId = user.Id
		roomUser.Nickname = user.Username
		err = utils.InsertRecord(roomUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"room_uuid": room.Uuid,
	})
}

// ShowRoomHandler ...
func ShowRoomHandler(c *gin.Context) {
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
