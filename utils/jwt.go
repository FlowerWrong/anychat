package utils

import (
	"time"

	"github.com/FlowerWrong/anychat/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// doc https://godoc.org/github.com/dgrijalva/jwt-go

// CustomJWTClaims ...
type CustomJWTClaims struct {
	Username string `json:"username"`
	ID       int64  `json:"id"`
	UUID     string `json:"uuid"`
	jwt.StandardClaims
}

// GenerateToken ...
func GenerateToken(user *models.User) (string, error) {
	jwtSigningKey := []byte(viper.GetString("jwt_key"))

	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := CustomJWTClaims{
		user.Username,
		user.Id,
		user.Uuid,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSigningKey)
}

// ParseToken ...
func ParseToken(tokenStr string) (*CustomJWTClaims, error) {
	jwtSigningKey := []byte(viper.GetString("jwt_key"))
	token, err := jwt.ParseWithClaims(tokenStr, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSigningKey, nil
	})

	if claims, ok := token.Claims.(*CustomJWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
