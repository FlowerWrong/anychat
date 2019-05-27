package chat

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/util"
)

// PerformSingleChat ...
func PerformSingleChat(req Req, c *Client) (err error) {
	var singleChatCmd SingleChatCmd
	err = json.Unmarshal(req.Data, &singleChatCmd)
	if err != nil {
		return err
	}

	from, err := services.FindUserByUuid(singleChatCmd.From)
	if err != nil {
		return err
	}
	to, err := services.FindUserByUuid(singleChatCmd.To)
	if err != nil {
		return err
	}

	log.Println(from.Id, to.Id)

	chatMsg := new(models.ChatMessage)
	chatMsg.From = from.Id
	chatMsg.To = to.Id
	chatMsg.Uuid = util.UUID()
	chatMsg.Ack = req.Ack
	chatMsg.Content = singleChatCmd.Msg
	chatMsg.CreatedAt = time.Unix(0, singleChatCmd.CreatedAt)
	affected, err := db.Engine().Insert(chatMsg)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("affected not 1")
	}

	// ack response
	err = c.sendAckRes(req.Ack, WS_SINGLE_CHAT)
	if err != nil {
		log.Println(err)
		return err
	}

	// check to is online or not
	toClient, err := c.hub.FindClientByUserID(to.Id)
	if err != nil {
		log.Println(err)
		// offline

		// email and sms notification TODO
	} else {
		// online
		data, err := buildRes(req.Cmd, chatMsg.Uuid, SingleChatRes{UUID: chatMsg.Uuid, From: singleChatCmd.From, To: singleChatCmd.To, Msg: singleChatCmd.Msg, CreatedAt: singleChatCmd.CreatedAt})
		if err != nil {
			return err
		}
		toClient.send <- data
	}
	return nil
}
