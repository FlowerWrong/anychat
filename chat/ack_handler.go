package chat

import (
	"encoding/json"
	"time"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
)

func (c *Client) sendAckRes(ack string, action int32) error {
	data, err := buildRes(WS_ACK, ack, Ack{Action: action})
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}

// PerformAck ...
func PerformAck(req Req, c *Client) (err error) {
	var ack Ack
	err = json.Unmarshal(req.Data, &ack)
	if err != nil {
		return err
	}
	switch ack.Action {
	case WS_SINGLE_CHAT:
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
	}
	return nil
}
