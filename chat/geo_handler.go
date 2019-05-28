package chat

import (
	"encoding/json"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/utils"
)

// PerformGeo ...
func PerformGeo(req Req, c *Client) (err error) {
	var geoCmd GeoCmd
	err = json.Unmarshal(req.Data, &geoCmd)
	if err != nil {
		return err
	}

	updateLoginLog := new(models.LoginLog)
	updateLoginLog.Latitude = geoCmd.Latitude
	updateLoginLog.Longitude = geoCmd.Longitude
	err = utils.UpdateRecord(c.loginLogID, updateLoginLog)
	if err != nil {
		return err
	}

	geoRes := Res{Base: Base{Ack: req.Ack, Cmd: req.Cmd}, Data: json.RawMessage([]byte{})}
	data, err := json.Marshal(geoRes)
	if err != nil {
		return err
	}
	c.send <- data
	return nil
}
