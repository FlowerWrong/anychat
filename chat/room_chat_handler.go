package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/FlowerWrong/util"
)

// PerformRoomChat ...
func PerformRoomChat(req Req, c *Client) (err error) {
	var roomChatCmd RoomChatCmd
	err = json.Unmarshal(req.Data, &roomChatCmd)
	if err != nil {
		return err
	}

	fromUser, err := services.FindUserByUUID(roomChatCmd.From)
	if err != nil {
		return err
	}
	toRoom, err := services.FindRoomByUUID(roomChatCmd.To)
	if err != nil {
		return err
	}

	log.Println(fromUser.Id, toRoom.Id)

	chatMsg := new(models.RoomMessage)
	chatMsg.From = fromUser.Id
	chatMsg.RoomId = toRoom.Id
	chatMsg.Uuid = util.UUID()
	chatMsg.Content = roomChatCmd.Msg
	chatMsg.CreatedAt = time.Unix(0, roomChatCmd.CreatedAt)
	err = utils.InsertRecord(chatMsg)
	if err != nil {
		return err
	}

	roomUsers, err := services.FindRoomUserListByRoomID(toRoom.Id)
	if err != nil {
		return err
	}
	for _, ru := range roomUsers {
		urm := new(models.UserRoomMessage)
		urm.Uuid = util.UUID()
		urm.UserId = ru.UserId
		urm.RoomMessageId = chatMsg.Id
		urm.CreatedAt = chatMsg.CreatedAt
		if ru.UserId == c.userID {
			urm.ReadAt = time.Now()
		}
		err = utils.InsertRecord(urm)
		if err != nil {
			return err
		}

		if ru.UserId != c.userID {
			// check to is online or not
			toClient, err := c.hub.FindClientByUserID(ru.UserId)
			if err != nil {
				log.Println(err)
				// offline

				// email and sms notification TODO
			} else {
				// online
				data, err := buildRes(req.Cmd, urm.Uuid, RoomChatRes{UUID: chatMsg.Uuid, From: roomChatCmd.From, To: roomChatCmd.To, Msg: roomChatCmd.Msg, CreatedAt: roomChatCmd.CreatedAt})
				if err != nil {
					return err
				}
				toClient.send <- data
			}
		}
	}

	// ack response
	err = c.sendAckRes(req.Ack, TypeRoomChat)
	if err != nil {
		return err
	}
	return nil
}
