package services

import (
	"errors"

	"github.com/FlowerWrong/new_chat/venus/db"
	"github.com/FlowerWrong/new_chat/venus/models"
)

func FindUserByUuid(uuid string) (*models.User, error) {
	var user models.User
	has, err := db.Engine().Where("uuid = ?", uuid).Get(&user)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("record not found")
	}
	return &user, nil
}
