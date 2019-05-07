package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	_ "./docs"
	"github.com/FlowerWrong/new_chat/venus/actions"
	"github.com/FlowerWrong/new_chat/venus/chat"
	"github.com/FlowerWrong/new_chat/venus/config"
	"github.com/FlowerWrong/new_chat/venus/db"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title one chat Swagger API
// @version 1.0
// @description one chat api server with websocket
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
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

	hub := chat.NewHub()
	go hub.Run()

	app := gin.Default()

	// swagger middleware
	swaggerConfig := &ginSwagger.Config{
		URL: "http://localhost:8080/swagger/doc.json", // The url pointing to API definition
	}
	app.GET("/swagger/*any", ginSwagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler))

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
