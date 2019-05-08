package chat

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/FlowerWrong/new_chat/venus/db"
	"github.com/FlowerWrong/new_chat/venus/models"
	"github.com/FlowerWrong/new_chat/venus/services"
	"github.com/FlowerWrong/new_chat/venus/utils"
	"github.com/FlowerWrong/util"
	"github.com/gorilla/websocket"
	"github.com/mssola/user_agent"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	realIP    string
	companyID int64
	appID     int64
	userID    int64
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})
	for {
		msgType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		switch msgType {
		case websocket.TextMessage:
			log.Println("TextMessage", string(message))
			var req Req
			err = json.Unmarshal(message, &req)
			if err != nil {
				log.Println(err)
				break
			}
			switch req.Cmd {
			case WS_LOGIN:
				var loginCmd LoginCmd
				err = json.Unmarshal(req.Body, &loginCmd)
				if err != nil {
					log.Println(err)
					break
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

				affected, err := db.Engine().Insert(user)
				if err != nil {
					log.Println(err)
					break
				}
				if affected != 1 {
					log.Println("insert failed", affected)
					break
				}
				c.userID = user.Id // 设置client user id

				// check domain
				// sql := "select * from apps where ? <@ domains limit 1"
				// domains := []string{loginCmd.Domain}
				// var apps []models.App
				// err = db.Engine().SQL(sql, pq.Array(domains)).Find(&apps)
				// if err != nil {
				// 	log.Println(err)
				// 	break
				// }
				// if len(apps) == 0 {
				// 	log.Println("Can not find company app")
				// 	break
				// }
				// app := apps[0]

				app := new(models.App)
				_, err = db.Engine().Where("token = ?", loginCmd.Token).Get(app)
				if err != nil {
					log.Println(err)
					break
				}
				c.companyID = app.CompanyId
				c.appID = app.Id

				// 选择对象
				users := make([]models.User, 0)
				err = db.Engine().Where("role = 'member' and company_id = ? and app_id = ?", c.companyID, c.appID).Find(&users) // TODO online
				if err != nil {
					log.Println(err)
					break
				}
				u := users[rand.Intn(len(users))] // TODO

				loginRes := LoginRes{UserID: user.Uuid, ChatID: u.Uuid}
				data, err := json.Marshal(loginRes)
				if err != nil {
					log.Println(err)
					break
				}

				c.send <- data
			case WS_LOGOUT:
				// TODO
			case WS_GEO:
				var geoCmd GeoCmd
				err = json.Unmarshal(req.Body, &geoCmd)
				if err != nil {
					log.Println(err)
					break
				}

				user := new(models.User)
				user.Latitude = geoCmd.Latitude
				user.Longitude = geoCmd.Longitude
				_, err = db.Engine().Id(c.userID).Cols("latitude", "longitude").Update(user)
				if err != nil {
					log.Println(err)
					break
				}
			case WS_LAN_IP:
				var lanIPCmd LanIPCmd
				err = json.Unmarshal(req.Body, &lanIPCmd)
				if err != nil {
					log.Println(err)
					break
				}

				user := new(models.User)
				user.LanIp = lanIPCmd.LanIP
				_, err = db.Engine().Id(c.userID).Cols("lan_ip").Update(user)
				if err != nil {
					log.Println(err)
					break
				}
			case WS_SINGLE_CHAT:
				var singleChatCmd SingleChatCmd
				err = json.Unmarshal(req.Body, &singleChatCmd)
				if err != nil {
					log.Println(err)
					break
				}

				from, err := services.FindUserByUuid(singleChatCmd.From)
				if err != nil {
					log.Println(err)
					break
				}
				to, err := services.FindUserByUuid(singleChatCmd.To)
				if err != nil {
					log.Println(err)
					break
				}

				log.Println(from.Ip, to.Ip)

				chatMsg := new(models.ChatMessage)
				chatMsg.From = from.Id
				chatMsg.To = to.Id
				chatMsg.Uuid = util.UUID()
				chatMsg.Ack = req.Ack
				chatMsg.Content = singleChatCmd.Msg
				affected, err := db.Engine().Insert(chatMsg)
				if err != nil {
					log.Println(err)
					break
				}
				if affected != 1 {
					log.Println("insert failed", affected)
					break
				}

				// check to is online or not
				toClient, err := c.hub.FindClientByUserId(to.Id)
				if err != nil {
					log.Println(err)
					// offline

					// email and sms notification TODO
				} else {
					// online
					singleChatRes := SingleChatRes{UUID: chatMsg.Uuid, Cmd: req.Cmd, Ack: req.Ack, From: singleChatCmd.From, To: singleChatCmd.To, Msg: singleChatCmd.Msg}
					data, err := json.Marshal(singleChatRes)
					if err != nil {
						log.Println(err)
						break
					}
					toClient.send <- data

					// 标记已读
					chatMsg.ReadAt = time.Now()
					affected, err = db.Engine().Id(chatMsg.Id).Cols("read_at").Update(&chatMsg)
					if err != nil {
						log.Println(err)
						break
					}
					if affected != 1 {
						log.Println("update failed", affected)
						break
					}
				}
			}
		case websocket.BinaryMessage:
			log.Println("BinaryMessage")
		case websocket.CloseMessage:
			log.Println("CloseMessage")
		case websocket.PingMessage:
			log.Println("PingMessage")
		case websocket.PongMessage:
			log.Println("PongMessage")
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// HandleWs handles websocket requests from the peer.
func HandleWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	realIP := utils.RealIP(r)

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, maxMessageSize), realIP: realIP}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
