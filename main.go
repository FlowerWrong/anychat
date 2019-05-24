package main

import (
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

	err := config.Setup()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server launch in", config.ENV)

	dbVersion, err := db.Engine().SqlTemplateClient("version.tpl").Query().List()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dbVersion[0]["version"])

	err = db.Redis().SetNX("ping", "pong", 10*time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}
	pong, err := db.Redis().Get("ping").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("redis", pong)

	hub := chat.NewHub()
	go hub.Run()

	app := gin.Default()

	app.LoadHTMLGlob("views/*")
	app.GET("/", actions.HomeHandler)

	v1 := app.Group("/api/v1")
	{
		v1.GET("/ping", actions.PingHandler)
		v1.POST("/upload", actions.UploadHandler)
	}
	app.GET("/ws", func(c *gin.Context) {
		actions.WsHandler(hub, c.Writer, c.Request)
	})
	err = app.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
