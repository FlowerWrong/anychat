package chat

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/FlowerWrong/anychat/db"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/gorilla/websocket"
	"github.com/nats-io/go-nats"
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
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	stop chan time.Time

	// Buffered channel of outbound messages.
	send chan []byte

	realIP    string
	companyID int64
	appID     int64
	userID    int64
}

func (c *Client) logical(message []byte) error {
	var req Req
	err := json.Unmarshal(message, &req)
	if err != nil {
		return err
	}

	switch req.Cmd {
	case WS_LOGIN:
		err = PerformLogin(req, c)
		if err != nil {
			return err
		}
	case WS_LOGOUT:
		// TODO
	case WS_RE_CONN:
		// 掉线重连 TODO
	case WS_GEO:
		err = PerformGeo(req, c)
		if err != nil {
			return err
		}
	case WS_LAN_IP:
	case WS_SINGLE_CHAT:
		err = PerformSingleChat(req, c)
		if err != nil {
			return err
		}
	case WS_PING:
		err = PerformPing(req, c)
		if err != nil {
			return err
		}
	}
	return nil
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.stop <- time.Now()
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
			err = c.logical(message)
			if err != nil {
				log.Println(err)
				break
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

func (c *Client) subPump() {
	defer func() {
		close(c.stop)
		_ = c.conn.Close()
	}()

	subj := "topic/" + strconv.FormatInt(c.userID, 10)
	db.MQClient().Subscribe(subj, func(msg *nats.Msg) {
		log.Printf("Received on [%s] Pid[%d]: '%s'", msg.Subject, os.Getpid(), string(msg.Data))

		// 逻辑
		err := c.logical(msg.Data)
		if err != nil {
			log.Println(err)
			c.stop <- time.Now()
		}
	})
	db.MQClient().Flush()
	<-c.stop
}

// HandleWs handles websocket requests from the peer.
func HandleWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	realIP := utils.RealIP(r)

	client := &Client{hub: hub, conn: conn, stop: make(chan time.Time), send: make(chan []byte, maxMessageSize), realIP: realIP}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
