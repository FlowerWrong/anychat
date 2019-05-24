package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler ... TODO plugin, eg database, redis, nats check
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
