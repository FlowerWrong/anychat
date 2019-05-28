package services

import (
	"errors"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
)

// FindChatMessageByUUID ...
func FindChatMessageByUUID(uuid string) (*models.ChatMessage, error) {
	var cm models.ChatMessage
	has, err := db.Engine().Where("uuid = ?", uuid).Get(&cm)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("record not found")
	}
	return &cm, nil
}

// FindUserRoomMessageByUUID ...
func FindUserRoomMessageByUUID(uuid string) (*models.UserRoomMessage, error) {
	var urm models.UserRoomMessage
	has, err := db.Engine().Where("uuid = ?", uuid).Get(&urm)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("record not found")
	}
	return &urm, nil
}
