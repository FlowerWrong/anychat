package chat

import (
	"encoding/json"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
)

// PerformLANIP ...
func PerformLANIP(req Req, c *Client) (err error) {
	var lanIPCmd LanIPCmd
	err = json.Unmarshal(req.Data, &lanIPCmd)
	if err != nil {
		return err
	}

	user := new(models.User)
	user.LanIp = lanIPCmd.LanIP
	_, err = db.Engine().Id(c.userID).Cols("lan_ip").Update(user)
	if err != nil {
		return err
	}
	lanIPRes := Res{Base: Base{Ack: req.Ack, Cmd: req.Cmd}, Data: json.RawMessage([]byte{})}
	data, err := json.Marshal(lanIPRes)
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}
