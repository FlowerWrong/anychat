package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler ...
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tpl", gin.H{
		"title": "Main website",
	})
}
