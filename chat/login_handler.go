package chat

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/models"
	"github.com/FlowerWrong/util"
	"github.com/mssola/user_agent"
)

// PerformLogin ...
func PerformLogin(req Req, c *Client) (err error) {
	var loginCmd LoginCmd
	err = json.Unmarshal(req.Body, &loginCmd)
	if err != nil {
		return err
	}
	ua := user_agent.New(loginCmd.UserAgent)
	user := new(models.User)
	browserName, browserVersion := ua.Browser()
	user.Browser = browserName + ":" + browserVersion
	user.Os = ua.OS()
	user.Ip = c.realIP

	user.Uuid = util.UUID()
	user.FirstLoginAt = time.Now()
	user.LastActiveAt = time.Now()

	// 可能存在，依赖于用户设置
	user.Username = loginCmd.Username
	user.Email = loginCmd.Email
	user.Mobile = loginCmd.Mobile

	affected, err := db.Engine().Insert(user)
	if err != nil {
		return err
	}
	if affected != 1 {
		return errors.New("affected not 1")
	}
	c.userID = user.Id // 设置client user id

	loginRes := LoginRes{UserID: user.Uuid, ChatID: ""}
	data, err := json.Marshal(loginRes)
	if err != nil {
		return err
	}

	c.send <- data
	return nil
}
