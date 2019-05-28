package chat

import (
	"encoding/json"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/utils"
)

// PerformLANIP ...
func PerformLANIP(req Req, c *Client) (err error) {
	var lanIPCmd LanIPCmd
	err = json.Unmarshal(req.Data, &lanIPCmd)
	if err != nil {
		return err
	}

	updateLoginLog := new(models.LoginLog)
	updateLoginLog.LanIp = lanIPCmd.LanIP
	err = utils.UpdateRecord(c.loginLogID, updateLoginLog)
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
