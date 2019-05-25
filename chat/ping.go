package chat

import (
	"encoding/json"
	"time"

	"github.com/FlowerWrong/anychat/utils"

	"github.com/FlowerWrong/util"
)

// doc https://docs.anycable.io/#/action_cable_protocol

func (c *Client) sendPing() {
	c.send <- newPingMessage()
	c.addPing()
}

func (c *Client) addPing() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}

	c.pingTimer = time.AfterFunc(appLayerpingInterval, c.sendPing)
}

func newPingMessage() []byte {
	raw, err := utils.RawMsg(PingCmd{PingAt: time.Now().Unix()})
	if err != nil {
		panic("Failed to build raw JSON 😲")
	}
	pingCmd := Res{Base: Base{Ack: util.UUID(), Cmd: WS_SERVER_PING}, Data: raw}
	jsonStr, err := json.Marshal(&pingCmd)
	if err != nil {
		panic("Failed to build ping JSON 😲")
	}
	return jsonStr
}
