package main

import (
	. "gen/config"
	"gen/controller"
	"gen/router"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"syscall"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	// 加载路由
	router.Route(engine)
	// 启动服务器，grace restart
	server := endless.NewServer(":"+Conf.String("APP_PORT"), engine)
	// 注册程序终止信号
	var signals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	}
	for _, signal := range signals {
		_ = server.RegisterSignalHook(endless.PRE_SIGNAL, signal, func() {
			controller.ArticleController.StopTheWorld()
		})
	}
	log.Println("server started success")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("server start failed, error:%s", err.Error())
	}
}
