package middlewares

import (
	"net/http"
	"strings"

	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuth ...
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.DefaultQuery("token", "")
		if token == "" {
			token = c.Request.Header.Get("Authorization")
			if s := strings.Split(token, " "); len(s) == 2 {
				token = s[1]
			}
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		user, err := services.FindUserByUUID(claims.UUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Set("currentUser", user)
		c.Next()
	}
}
