package chat

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/FlowerWrong/god/db"
	"github.com/FlowerWrong/new_chat/venus/models"
	"github.com/FlowerWrong/new_chat/venus/utils"
	"github.com/FlowerWrong/util"
	"github.com/gorilla/websocket"
	"github.com/lib/pq"
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
			log.Println("TextMessage", message)
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
				user.LanIp = loginCmd.LanIP
				user.Latitude = loginCmd.Latitude
				user.Longitude = loginCmd.Longitude

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

				// 选择对象
				users := make([]models.User, 0)
				err = db.Engine().Where("role > member and company_id = ? and app_id = ?", c.companyID, c.appID).Find(&users)
				if err != nil {
					log.Println(err)
					break
				}
				u := users[rand.Intn(len(users))] // TODO

				loginRes := LoginRes{Uuid: user.Uuid, ChatID: u.Uuid}
				data, err := json.Marshal(loginRes)
				if err != nil {
					log.Println(err)
					break
				}

				c.hub.broadcast <- data
			case WS_LOGOUT:
			case WS_SINGLE_CHAT:
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

	domain := utils.Host(r)
	sql := "select * from apps where ? <@ domains limit 1"
	domains := []string{domain}
	var apps []models.App
	err = db.Engine().SQL(sql, pq.Array(domains)).Find(&apps)
	if err != nil {
		log.Println(err)
		return
	}
	if len(apps) == 0 {
		log.Println("Can not find company app")
		return
	}
	app := apps[0]

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, maxMessageSize), realIP: realIP, companyID: app.CompanyId, appID: app.Id}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
