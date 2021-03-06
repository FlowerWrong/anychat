package chat

import (
	"encoding/json"
	"time"

	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/anychat/services"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/mssola/user_agent"
)

// PerformLogin ...
func (c *Client) PerformLogin(req Req) (err error) {
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

	user, err := services.FindUserByUUID(claims.UUID)
	if err != nil {
		return err
	}

	updateUser := new(models.User)
	updateUser.FirstLoginAt = time.Now()
	updateUser.LastActiveAt = time.Now()
	err = utils.UpdateRecord(user.Id, updateUser)
	if err != nil {
		return err
	}

	newLoginLog := new(models.LoginLog)
	newLoginLog.UserId = user.Id
	browserName, browserVersion := ua.Browser()
	newLoginLog.Browser = browserName + ":" + browserVersion
	newLoginLog.Os = ua.OS()
	newLoginLog.Ip = c.realIP
	err = utils.InsertRecord(newLoginLog)
	if err != nil {
		return err
	}

	c.userID = user.Id // 设置client user id
	c.userUUID = user.Uuid
	c.loginLogID = newLoginLog.Id

	data, err := buildRes(req.Cmd, req.Ack, LoginRes{UserID: user.Uuid})
	if err != nil {
		return err
	}

	c.send <- data
	c.updateLogined()
	return nil
}
