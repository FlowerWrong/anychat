package chat

import (
	"encoding/json"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
)

// PerformGeo ...
func PerformGeo(req Req, c *Client) (err error) {
	var geoCmd GeoCmd
	err = json.Unmarshal(req.Body, &geoCmd)
	if err != nil {
		return err
	}

	user := new(models.User)
	user.Latitude = geoCmd.Latitude
	user.Longitude = geoCmd.Longitude
	_, err = db.Engine().Id(c.userID).Cols("latitude", "longitude").Update(user)
	if err != nil {
		return err
	}
	return nil
}
