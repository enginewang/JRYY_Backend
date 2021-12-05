package main

import (
	"JRYY/db"
	"JRYY/server"
	"log"
)

func main() {
	//r := gin.Default()
	//r.GET("/", func(context *gin.Context) {
	//	context.String(http.StatusOK, "Hello Gin!")
	//})
	err := db.InitGlobalDatabase("127.0.0.1", "jryy")
	if err != nil {
		log.Panic(err)
	}
	//err = db.InitRedisServer()
	//if err != nil {
	//	log.Panic(err)
	//}
	r := server.NewServer(":7439")
	err = r.Init()
	if err != nil {
		log.Panic(err)
	}
	err = r.StartServer()
	if err != nil {
		log.Panic(err)
	}
}