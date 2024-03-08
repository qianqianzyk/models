package main

import (
	"github.com/gin-gonic/gin"
	"live-chat/app/midwares"
	"live-chat/config/database"
	"live-chat/config/router"
	"log"
	"math/rand"
	"time"
)

func main() {
	database.Init()
	//initFile.InitFile()
	r := gin.Default()
	r.Use(midwares.Cors())
	r.Use(midwares.Corss)
	r.Use(midwares.RateLimitMiddleware(time.Second, 100, 100))
	r.Use(midwares.ErrHandler())
	r.NoMethod(midwares.HandleNotFound)
	r.NoRoute(midwares.HandleNotFound)
	rand.Seed(time.Now().UnixNano())
	router.Init(r)
	err := r.Run()
	if err != nil {
		log.Fatal("ServerStartFailed", err)
	}
}
