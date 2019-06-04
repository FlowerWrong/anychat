package actions

import (
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

// TODO 事务

// CreateRoomHandler ...
func CreateRoomHandler(c *gin.Context) {
	var roomForm roomVM
	if err := c.ShouldBindJSON(&roomForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUserI, exists := c.Get("currentUser")
	if exists == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current user not found"})
		return
	}
	currentUser := currentUserI.(*models.User)

	room := new(models.Room)
	room.Uuid = util.UUID()
	room.Name = roomForm.Name
	room.Intro = roomForm.Intro
	room.CreatorId = currentUser.Id
	room.Logo = roomForm.Logo
	err := utils.InsertRecord(room)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 添加自己
	if !utils.Contains(roomForm.UserUUIDs, currentUser.Uuid) {
		roomForm.UserUUIDs = append(roomForm.UserUUIDs, currentUser.Uuid)
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

	c.JSON(http.StatusOK, room)
}

// ShowRoomHandler ...
func ShowRoomHandler(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if exists == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current user not found"})
		return
	}
	currentUser := currentUserI.(*models.User)

	uuids, err := services.FindMyRoomUUIDList(currentUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	uuid := c.Param("uuid")
	if !utils.Contains(uuids, uuid) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "403",
		})
		return
	}

	room, err := services.FindRoomByUUID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, room)
}

// IndexRoomHandler ...
func IndexRoomHandler(c *gin.Context) {
	currentUserI, exists := c.Get("currentUser")
	if exists == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "current user not found"})
		return
	}
	currentUser := currentUserI.(*models.User)

	rooms, err := services.FindMyRoomList(currentUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rooms)
}
