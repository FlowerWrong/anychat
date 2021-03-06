package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/FlowerWrong/anychat/chat"
	"github.com/FlowerWrong/anychat/utils"
	"github.com/FlowerWrong/util"
	"github.com/gorilla/websocket"
)

// Session ...
type Session struct {
	conn   *websocket.Conn
	send   chan []byte
	mu     sync.Mutex
	smap   sync.Map
	stage  int32
	userID string
	token  string
}

func (s *Session) updateStage(stage int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stage = stage
}

type loginVM struct {
	Username string `form:"username" json:"username" xml:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

func (s *Session) httpLogin(username, password string) error {
	loginFrom := loginVM{Username: username, Password: password}
	postLoginJSON, _ := json.Marshal(loginFrom)
	loginURL := "http://localhost:8080/api/v1/login"

	loginReq, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(postLoginJSON))
	loginReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	loginResp, err := client.Do(loginReq)
	if err != nil {
		log.Println(err)
		return err
	}
	defer loginResp.Body.Close()
	loginBody, _ := ioutil.ReadAll(loginResp.Body)
	var m map[string]interface{}
	err = json.Unmarshal(loginBody, &m)
	if err != nil {
		log.Println(err)
		return err
	}
	s.token = m["token"].(string)
	return nil
}

func (s *Session) sendChatMsg(chatToUUID, msg string) error {
	raw, err := utils.RawMsg(chat.SingleChatCmd{From: s.userID, To: chatToUUID, Msg: msg, CreatedAt: time.Now().UnixNano()})
	if err != nil {
		log.Println(err)
		return err
	}
	chatReq := chat.Req{Base: chat.Base{Ack: "single_chat", Cmd: chat.TypeSingleChat}, Data: raw}
	chatJSON, err := json.Marshal(chatReq)
	if err != nil {
		log.Println(err)
		return err
	}
	s.send <- chatJSON
	return nil
}

func (s *Session) sendRoomChatMsg(chatToRoomUUID, msg string) error {
	raw, err := utils.RawMsg(chat.RoomChatCmd{From: s.userID, To: chatToRoomUUID, Msg: msg, CreatedAt: time.Now().UnixNano()})
	if err != nil {
		log.Println(err)
		return err
	}
	chatReq := chat.Req{Base: chat.Base{Ack: "room_chat", Cmd: chat.TypeRoomChat}, Data: raw}
	chatJSON, err := json.Marshal(chatReq)
	if err != nil {
		log.Println(err)
		return err
	}
	s.send <- chatJSON
	return nil
}

func (s *Session) sendAckRes(ack, action string) error {
	raw, err := utils.RawMsg(chat.Ack{Action: action})
	if err != nil {
		log.Println(err)
		return err
	}
	req := chat.Req{Base: chat.Base{Ack: ack, Cmd: chat.TypeAck}, Data: raw}
	data, err := json.Marshal(req)
	if err != nil {
		log.Println(err)
		return err
	}
	s.send <- data
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	// flag
	username := flag.String("username", "", "login username")
	password := flag.String("password", "", "login password")
	chatToUUID := flag.String("to", "", "chat to someone's uuid")
	chatToRoomUUID := flag.String("room", "", "chat to room's uuid")
	flag.Parse()

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/anychat", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	session := &Session{conn: conn, send: make(chan []byte, chat.MaxMessageSize), stage: 1}
	defer func() {
		close(session.send)
		_ = session.conn.Close()
	}()

	err = session.httpLogin(*username, *password)
	if err != nil {
		log.Fatal("login:", err)
	}

	wgw := new(util.WaitGroupWrapper)
	wgw.Wrap(func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}

			var res chat.Res
			err = json.Unmarshal(message, &res)
			if err != nil {
				log.Println(err)
				break
			}

			var m map[string]interface{}
			err = json.Unmarshal(res.Data, &m)
			if err != nil {
				log.Println(err)
				break
			}
			switch res.Cmd {
			case chat.TypeLogin:
				// log.Println(m["user_id"])
				session.userID = m["user_id"].(string)
				session.updateStage(3)
			case chat.TypePing:
				// log.Println(m["ping_at"])
			case chat.TypeSingleChat:
				log.Println("chat", m["from"], "say", m["msg"], "to you")
				session.sendAckRes(res.Ack, chat.TypeSingleChat)
			case chat.TypeRoomChat:
				log.Println("room chat", m["from"], "say", m["msg"], "to you")
				session.sendAckRes(res.Ack, chat.TypeRoomChat)
			case chat.TypeAck:
				log.Println(m["action"].(string))
				switch m["action"] {
				case chat.TypeSingleChat:
					log.Println("single chat ack")
				case chat.TypeRoomChat:
					log.Println("root chat ack")
				}
			}
		}
	})

	wgw.Wrap(func() {
		for {
			select {
			case message, ok := <-session.send:
				_ = session.conn.SetWriteDeadline(time.Now().Add(chat.WriteWait))
				if !ok {
					_ = session.conn.WriteMessage(websocket.CloseMessage, []byte{})
					break
				}

				w, err := session.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Println(err)
					break
				}
				_, err = w.Write(message)
				if err := w.Close(); err != nil {
					log.Println(err)
					break
				}
			}
		}
	})

	wgw.Wrap(func() {
		for {
			switch session.stage {
			case 1:
				raw, err := utils.RawMsg(chat.LoginCmd{
					UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
					Domain:    "github.com",
					Token:     session.token,
				})
				if err != nil {
					log.Println(err)
					break
				}
				loginReq := chat.Req{Base: chat.Base{Ack: "login", Cmd: chat.TypeLogin}, Data: raw}
				loginJSON, err := json.Marshal(loginReq)
				if err != nil {
					log.Println(err)
					break
				}
				session.send <- loginJSON
				session.updateStage(2)
			case 3:
				err = session.sendChatMsg(*chatToUUID, "Hello")
				if err != nil {
					break
				}
				session.updateStage(5)
			case 5:
				err = session.sendRoomChatMsg(*chatToRoomUUID, "Hi!, every one.")
				if err != nil {
					break
				}
				session.updateStage(6)
			case 7:
			case 9:
				break
			}
		}
	})

	reader := bufio.NewReader(os.Stdin)
	wgw.Wrap(func() {
		for {
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			msg := strings.Replace(text, "\n", "", -1)
			log.Println(msg)

			err = session.sendChatMsg(*chatToUUID, msg)
			if err != nil {
				break
			}
		}
	})

	wgw.WaitGroup.Wait()
}
