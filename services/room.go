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

// FindMyRoomList 我的群聊列表
func FindMyRoomList(currentUser *models.User) ([]models.Room, error) {
	paramMap := map[string]interface{}{"user_id": currentUser.Id}
	var rooms []models.Room
	err := db.Engine().SqlTemplateClient("my_room_list.tpl", &paramMap).Find(&rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// FindMyRoomUUIDList 我的群聊uuid列表
func FindMyRoomUUIDList(currentUser *models.User) ([]string, error) {
	rooms, err := FindMyRoomList(currentUser)
	if err != nil {
		return nil, err
	}
	var uuids []string
	for _, room := range rooms {
		uuids = append(uuids, room.Uuid)
	}
	return uuids, nil
}

// FindMyCreatedRoomList 我创建的room列表
func FindMyCreatedRoomList(currentUser *models.User) ([]models.Room, error) {
	var rooms []models.Room
	err := db.Engine().Where("creator_id = ?", currentUser.Id).Find(&rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}
