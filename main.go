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
	"github.com/FlowerWrong/anychat/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
	app.Use(middlewares.RateLimit())

	app.LoadHTMLGlob("views/*")
	app.GET("/", actions.HomeHandler)

	v1 := app.Group("/api/v1")
	{
		v1.GET("/health", actions.HealthHandler)
		v1.POST("/login", actions.LoginHandler)
	}

	authGroup := app.Group("/api/v1", middlewares.JWTAuth())
	{
		authGroup.POST("/upload", actions.UploadHandler)

		authGroup.GET("/rooms", actions.IndexRoomHandler)
		authGroup.POST("/rooms", actions.CreateRoomHandler)
		authGroup.GET("/rooms/:uuid", actions.ShowRoomHandler)

		// 获取历史聊天记录
		authGroup.GET("/rooms/:uuid/messages", actions.RoomChatMsgHandler)
		authGroup.GET("/messages", actions.SingleChatMsgHandler)
	}

	app.GET("/anychat", func(c *gin.Context) {
		actions.WsHandler(hub, c.Writer, c.Request)
	})

	err = app.Run(viper.GetString("server_url"))
	if err != nil {
		log.Fatal(err)
	}
}
