package chat

import (
	"encoding/json"

	"github.com/FlowerWrong/anychat/utils"
	"github.com/FlowerWrong/util"
)

func (c *Client) sendWelcome() {
	if c.closed {
		return
	}
	c.send <- newWelcomeMessage()
}

func newWelcomeMessage() []byte {
	raw, err := utils.RawMsg(WelcomeCmd{Message: "hello"})
	if err != nil {
		panic("Failed to build raw JSON 😲")
	}
	welcomeCmd := Res{Base: Base{Ack: util.UUID(), Cmd: WS_WELCOME}, Data: raw}
	jsonStr, err := json.Marshal(&welcomeCmd)
	if err != nil {
		panic("Failed to build ping JSON 😲")
	}
	return jsonStr
}
