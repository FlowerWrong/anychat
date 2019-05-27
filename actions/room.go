package actions

import (
	"errors"
	"net/http"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
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
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

// ShowRoomHandler ...
func ShowRoomHandler(c *gin.Context) {
	id := c.Param("id")
	room := new(models.Room)
	has, err := db.Engine().Id(id).Get(room)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	if !has {
		c.JSON(http.StatusNotFound, gin.H{
			"message": errors.New("record not found"),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"room_id": room.Id,
	})
}
