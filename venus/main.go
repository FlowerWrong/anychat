package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/FlowerWrong/new_chat/venus/actions"
	"github.com/FlowerWrong/new_chat/venus/chat"
	"github.com/FlowerWrong/new_chat/venus/config"
	"github.com/FlowerWrong/new_chat/venus/db"
	"github.com/FlowerWrong/new_chat/venus/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := config.Setup()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(viper.Get("db_url"), config.ENV)

	dbVersion, err := db.Engine().SqlTemplateClient("version.tpl").Query().Json()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(dbVersion)

	err = db.Redis().SetNX("ping", "pong", 10*time.Second).Err()
	if err != nil {
		log.Fatal(err)
	}
	pong, err := db.Redis().Get("ping").Result()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(pong)

	var users []models.User
	err = db.Engine().Where("deleted_at is NULL").Limit(10, 0).Find(&users)
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range users {
		log.Println(user.Username)
	}

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
