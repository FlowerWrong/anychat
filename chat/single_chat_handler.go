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
	err = json.Unmarshal(req.Body, &singleChatCmd)
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

	log.Println(from.Ip, to.Ip)

	chatMsg := new(models.ChatMessage)
	chatMsg.From = from.Id
	chatMsg.To = to.Id
	chatMsg.Uuid = util.UUID()
	chatMsg.Ack = req.Ack
	chatMsg.Content = singleChatCmd.Msg
	affected, err := db.Engine().Insert(chatMsg)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("affected not 1")
	}

	// check to is online or not
	toClient, err := c.hub.FindClientByUserID(to.Id)
	if err != nil {
		log.Println(err)
		// offline

		// email and sms notification TODO
	} else {
		// online
		singleChatRes := SingleChatRes{UUID: chatMsg.Uuid, Cmd: req.Cmd, Ack: req.Ack, From: singleChatCmd.From, To: singleChatCmd.To, Msg: singleChatCmd.Msg}
		data, err := json.Marshal(singleChatRes)
		if err != nil {
			log.Println(err)
			return err
		}
		toClient.send <- data

		// 标记已读
		chatMsg.ReadAt = time.Now()
		affected, err = db.Engine().Id(chatMsg.Id).Cols("read_at").Update(&chatMsg)
		if err != nil {
			return err
		}
		if affected != 1 {
			return errors.New("affected not 1")
		}
	}
	return nil
}
