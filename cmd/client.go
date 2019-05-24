package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/FlowerWrong/util"
	"github.com/gorilla/websocket"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	wgw := new(util.WaitGroupWrapper)
	wgw.Wrap(func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	})

	loginJSON := []byte(`{
		"ack": "123abc",
		"cmd": 0,
		"body": {
			"userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36",
			"domain": "github.com",
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InlhbmciLCJpZCI6MTAxLCJ1dWlkIjoiOTQxYWI2MmUtYmI4MS00MjNjLWI2ZTctNzVmZmU3OTdiNjBhIiwiZXhwIjoxNTU5MzIzNTI3LCJpc3MiOiJ0ZXN0In0.cPr7uqT1Mecp2DbXamVrFpAfQGETr-7UXRjO-NFDw90"
		}
	}`)
	err = conn.WriteMessage(websocket.TextMessage, loginJSON)
	if err != nil {
		log.Println("write:", err)
		return
	}

	wgw.WaitGroup.Wait()
}
