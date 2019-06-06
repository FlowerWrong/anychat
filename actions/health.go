package actions

import (
	"net/http"

	"github.com/FlowerWrong/anychat/db"
	"github.com/gin-gonic/gin"
)

// HealthHandler ...
func HealthHandler(c *gin.Context) {
	var redisStatus string
	_, err := db.Redis().Ping().Result()
	if err != nil {
		redisStatus = err.Error()
	} else {
		redisStatus = "running"
	}

	var dbStatus string
	_, err = db.Engine().SqlTemplateClient("version.tpl").Query().List()
	if err != nil {
		dbStatus = err.Error()
	} else {
		dbStatus = "running"
	}

	c.JSON(http.StatusOK, gin.H{
		"app":      "running",
		"redis":    redisStatus,
		"database": dbStatus,
	})
}
