package chat

import (
	"encoding/json"
	"time"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
)

// PerformAck ...
func PerformAck(req Req, c *Client) (err error) {
	var ack Ack
	err = json.Unmarshal(req.Data, &ack)
	if err != nil {
		return err
	}
	switch ack.Action {
	case TypeSingleChat:
		uuid := req.Ack
		cm, err := services.FindChatMessageByUUID(uuid)
		if err != nil {
			return err
		}

		// 标记已读
		updateMsg := new(models.ChatMessage)
		updateMsg.ReadAt = time.Now()
		err = utils.UpdateRecord(cm.Id, updateMsg)
		if err != nil {
			return err
		}
	case TypeRoomChat:
		uuid := req.Ack
		urm, err := services.FindUserRoomMessageByUUID(uuid)
		if err != nil {
			return err
		}

		// 标记已读
		updateMsg := new(models.UserRoomMessage)
		updateMsg.ReadAt = time.Now()
		err = utils.UpdateRecord(urm.Id, updateMsg)
		if err != nil {
			return err
		}
	}
	return nil
}
