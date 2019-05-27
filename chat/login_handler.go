package chat

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"

	"github.com/mssola/user_agent"
)

// PerformLogin ...
func PerformLogin(req Req, c *Client) (err error) {
	var loginCmd LoginCmd
	err = json.Unmarshal(req.Data, &loginCmd)
	if err != nil {
		return err
	}
	ua := user_agent.New(loginCmd.UserAgent)

	claims, err := utils.ParseToken(loginCmd.Token)
	if err != nil {
		return err
	}

	user, err := services.FindUserByUuid(claims.UUID)
	if err != nil {
		return err
	}

	updateUser := new(models.User)
	browserName, browserVersion := ua.Browser()
	updateUser.Browser = browserName + ":" + browserVersion
	updateUser.Os = ua.OS()
	updateUser.Ip = c.realIP
	updateUser.FirstLoginAt = time.Now()
	updateUser.LastActiveAt = time.Now()

	affected, err := db.Engine().Id(user.Id).Update(updateUser)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("affected not 1")
	}
	c.userID = user.Id // 设置client user id
	c.userUUID = user.Uuid

	data, err := buildRes(req.Cmd, req.Ack, LoginRes{UserID: user.Uuid})
	if err != nil {
		return err
	}

	c.send <- data
	return nil
}
