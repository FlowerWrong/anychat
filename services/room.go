package services

import (
	"errors"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
)

// FindRoomByUUID ...
func FindRoomByUUID(uuid string) (*models.Room, error) {
	var r models.Room
	has, err := db.Engine().Where("uuid = ?", uuid).Get(&r)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("record not found")
	}
	return &r, nil
}

// FindRoomUserListByRoomID ...
func FindRoomUserListByRoomID(id int64) ([]models.RoomUser, error) {
	var ru []models.RoomUser
	err := db.Engine().Where("room_id = ?", id).Find(&ru)
	if err != nil {
		return nil, err
	}
	return ru, nil
}
