package chat

import (
	"encoding/json"
	"log"
)

func (c *Client) sendAckRes(ack string) error {
	ackJSON := Res{Base: Base{Cmd: WS_ACK, Ack: ack}, Data: []byte{}}
	data, err := json.Marshal(ackJSON)
	if err != nil {
		log.Println(err)
		return err
	}
	c.send <- data
	return nil
}

// PerformAck ...
func PerformAck(req Req, c *Client) (err error) {
	return nil
}
