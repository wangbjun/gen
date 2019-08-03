package main

import (
	. "gen/config"
	"gen/router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// release mode
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	// 加载路由
	router.Route(engine)
	// 启动服务器，grace restart
	server := endless.NewServer(":"+Conf.String("APP_PORT"), engine)
	log.Println("server started success")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server start failed, error:%s", err.Error())
	}
}
