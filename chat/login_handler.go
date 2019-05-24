package chat

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
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
	user.CompanyId = c.companyID
	user.Role = "customer"
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

	app := new(models.App)
	_, err = db.Engine().Where("token = ?", loginCmd.Token).Get(app)
	if err != nil {
		return err
	}
	c.companyID = app.CompanyId
	c.appID = app.Id

	// 选择对象
	users := make([]models.User, 0)
	err = db.Engine().Where("role = 'member' and company_id = ? and app_id = ?", c.companyID, c.appID).Find(&users) // TODO online
	if err != nil {
		return err
	}
	if len(users) == 0 {
		return errors.New("no company users find")
	}
	onlineUsers := c.hub.FindOnlineUserList(&users)
	var selectedUser models.User
	if len(onlineUsers) == 0 {
		log.Println("no online users find")
		selectedUser = users[rand.Intn(len(users))]
	} else {
		selectedUser = *onlineUsers[rand.Intn(len(onlineUsers))]
	}

	loginRes := LoginRes{UserID: user.Uuid, ChatID: selectedUser.Uuid}
	data, err := json.Marshal(loginRes)
	if err != nil {
		return err
	}

	c.send <- data
	go c.subPump()
	return nil
}
