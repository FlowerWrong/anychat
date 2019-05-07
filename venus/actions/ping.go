package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PingHandler ...
// @Summary ping pong for test
// @Description req ping, res pong
// @Accept json
// @Produce json
// @Success 200 {string} string	"pong"
// @Router /ping [get]
func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
