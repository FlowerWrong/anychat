package actions

import (
	"net/http"
	"time"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/gin-gonic/gin"
)

type loginVM struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

// LoginHandler ...
// curl -v -X POST http://localhost:8080/api/v1/login -H 'content-type: application/json' -d '{ "username": "yang", "password": "123456" }'
func LoginHandler(c *gin.Context) {
	var login loginVM
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := services.FindUserByUsername(login.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = utils.CheckPassword(user.PasswordDigest, login.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password is wrong"})
		return
	}

	token, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	statusCmd := db.Redis().Set(user.Uuid, token, time.Duration(1000*60*60*24*7))
	if statusCmd.Err() != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": statusCmd.Err().Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
