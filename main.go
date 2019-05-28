package main

import (
	"flag"
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/FlowerWrong/anychat/actions"
	"github.com/FlowerWrong/anychat/chat"
	"github.com/FlowerWrong/anychat/config"
	"github.com/FlowerWrong/anychat/db"
	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	configFile := flag.String("config", "./config/settings.yml", "config file path")
	flag.Parse()

	err := config.Setup(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server launch in", config.APP_ENV)

	dbVersion, err := db.Engine().SqlTemplateClient("version.tpl").Query().List()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dbVersion[0]["version"])

	hub := chat.NewHub()
	go hub.Run()

	app := gin.Default()

	app.LoadHTMLGlob("views/*")
	app.GET("/", actions.HomeHandler)

	v1 := app.Group("/api/v1")
	{
		v1.GET("/health", actions.HealthHandler)
		v1.POST("/login", actions.LoginHandler)
		v1.POST("/upload", actions.UploadHandler)

		v1.POST("/rooms", actions.CreateRoomHandler)
		v1.GET("/rooms/:uuid", actions.ShowRoomHandler)

		// 获取历史聊天记录
		v1.GET("/rooms/:uuid/messages", actions.RoomChatMsgHandler)
		v1.GET("/messages", actions.SingleChatMsgHandler)
	}
	app.GET("/anychat", func(c *gin.Context) {
		actions.WsHandler(hub, c.Writer, c.Request)
	})
	err = app.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
