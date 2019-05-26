package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/FlowerWrong/anychat/actions"
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
	stage  int32
	userID string
	token  string
}

func (s *Session) updateStage(stage int32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.stage = stage
}

func (s *Session) login(username, password string) {
	loginFrom := actions.Login{Username: username, Password: password}
	postLoginJSON, _ := json.Marshal(loginFrom)
	log.Println(string(postLoginJSON))
	loginURL := "http://localhost:8080/api/v1/login"

	loginReq, err := http.NewRequest("POST", loginURL, bytes.NewBuffer(postLoginJSON))
	loginReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	loginResp, err := client.Do(loginReq)
	if err != nil {
		panic(err)
	}
	defer loginResp.Body.Close()
	loginBody, _ := ioutil.ReadAll(loginResp.Body)
	var m map[string]interface{}
	err = json.Unmarshal(loginBody, &m)
	if err != nil {
		log.Println(err)
		return
	}
	s.token = m["token"].(string)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	// flag
	username := flag.String("username", "", "login username")
	password := flag.String("password", "", "login password")
	chatToUUID := flag.String("to", "", "chat to someone's uuid")
	flag.Parse()

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	session := &Session{conn: conn, send: make(chan []byte, chat.MaxMessageSize), stage: 1}
	defer func() {
		close(session.send)
		_ = session.conn.Close()
	}()

	session.login(*username, *password)

	wgw := new(util.WaitGroupWrapper)
	wgw.Wrap(func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
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
				return
			}
			switch res.Cmd {
			case chat.WS_LOGIN:
				// log.Println(m["user_id"])
				session.userID = m["user_id"].(string)
				session.updateStage(3)
			case chat.WS_SERVER_PING:
				// log.Println(m["ping_at"])
			case chat.WS_SINGLE_CHAT:
				log.Println(m["from"], "say", m["msg"], "to you")
			}
		}
	})

	wgw.Wrap(func() {
		for {
			select {
			case message, ok := <-session.send:
				_ = session.conn.SetWriteDeadline(time.Now().Add(chat.WriteWait))
				if !ok {
					// The hub closed the channel.
					_ = session.conn.WriteMessage(websocket.CloseMessage, []byte{})
					return
				}

				w, err := session.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				_, _ = w.Write(message)

				// Add queued chat messages to the current websocket message.
				n := len(session.send)
				for i := 0; i < n; i++ {
					_, _ = w.Write(<-session.send)
				}

				if err := w.Close(); err != nil {
					return
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
					log.Panic(err)
				}
				loginReq := chat.Req{Base: chat.Base{Ack: "single_chat", Cmd: chat.WS_LOGIN}, Data: raw}
				loginJSON, err := json.Marshal(loginReq)
				if err != nil {
					log.Panic(err)
				}
				session.send <- loginJSON
				session.updateStage(2)
			case 3:
				log.Println(3)
				raw, err := utils.RawMsg(chat.SingleChatCmd{From: session.userID, To: *chatToUUID, Msg: "Hello"})
				if err != nil {
					log.Panic(err)
				}
				loginRes := chat.Req{Base: chat.Base{Ack: "single_chat", Cmd: chat.WS_SINGLE_CHAT}, Data: raw}
				chatJSON, err := json.Marshal(loginRes)
				if err != nil {
					log.Panic(err)
				}
				session.send <- chatJSON
				session.updateStage(4)
			case 5:
			case 7:
			case 9:
				break
			}
		}
	})

	reader := bufio.NewReader(os.Stdin)
	wgw.Wrap(func() {
		for {
			fmt.Print("-> ")
			text, _ := reader.ReadString('\n')
			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			if strings.Compare("hi", text) == 0 {
				fmt.Println("hello, Yourself")
			}
			log.Println(text)
		}
	})

	wgw.WaitGroup.Wait()
}
