package main

import (
	"gen/config"
	"gen/router"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(getMode())
	engine := gin.New()
	engine.Use(gin.Recovery())
	// 加载路由
	router.Route(engine)
	// 启动服务器
	log.Println("server started success")
	err := engine.Run(":" + config.Conf.Section("APP").Key("PORT").String())
	if err != nil {
		log.Fatalf("server start failed, error: %s", err.Error())
	}
}

func getMode() string {
	debug := config.Conf.Section("APP").Key("DEBUG").String()
	if debug == "true" {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}
