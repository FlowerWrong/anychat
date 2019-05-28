package chat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/FlowerWrong/anychat/utils"
	"github.com/gorilla/websocket"
)

const (
	// WriteWait time allowed to write a message to the peer.
	WriteWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// MaxMessageSize maximum message size allowed from peer.
	MaxMessageSize = 1024

	appLayerpingInterval = 3 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  MaxMessageSize,
	WriteBufferSize: MaxMessageSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{"anychat-v1-json"},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	closed    bool
	connected bool

	realIP   string
	userID   int64
	userUUID string

	pingTimer *time.Timer

	mu sync.Mutex

	logined bool
}

func (c *Client) updateLogined() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logined = !c.logined
}

func (c *Client) logical(message []byte) error {
	var req Req
	err := json.Unmarshal(message, &req)
	if err != nil {
		return err
	}

	if req.Cmd == TypeLogin {
		err = PerformLogin(req, c)
		if err != nil {
			return err
		}
	} else {
		if !c.logined {
			return errors.New("401 Unauthorized")
		}
		switch req.Cmd {
		case TypeReConn: // without ack response
			// 掉线重连 TODO
		case TypeGeo: // without ack response
			err = PerformGeo(req, c)
			if err != nil {
				return err
			}
		case TypeLanIP: // without ack response
			err = PerformLANIP(req, c)
			if err != nil {
				return err
			}
		case TypeSingleChat: // with ack response
			err = PerformSingleChat(req, c)
			if err != nil {
				return err
			}
		case TypeRoomChat: // with ack response
			err = PerformRoomChat(req, c)
			if err != nil {
				return err
			}
		case TypeAck: // without ack response
			err = PerformAck(req, c)
			if err != nil {
				return err
			}
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
		c.closed = true
		c.hub.unregister <- c
		if c.pingTimer != nil {
			c.pingTimer.Stop()
		}
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(MaxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})

	c.sendPing()

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
				c.sendDisconnectRes(err.Error(), false)
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
			_ = c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
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
			_ = c.conn.SetWriteDeadline(time.Now().Add(WriteWait))
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

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, MaxMessageSize), realIP: realIP, connected: true, closed: false, logined: false}
	client.hub.register <- client
	client.sendWelcome()

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
