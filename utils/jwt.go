package utils

import (
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

	claims := CustomJWTClaims{
		user.Username,
		user.Id,
		user.Uuid,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSigningKey)
}

// ParseToken ...
func ParseToken(tokenStr string) (*CustomJWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})

	if claims, ok := token.Claims.(*CustomJWTClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
