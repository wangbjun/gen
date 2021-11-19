package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gen/config"
	"gen/log"
	"gen/models"
	"gen/router"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "conf", "app.ini", "config file path")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(configFile)
	if err != nil {
		panic(fmt.Sprintf("load config failed, file: %s, error: %s", configFile, err))
	}

	// 初始化日志
	log.Init(cfg)
	defer func() {
		if err := log.Logger.Sync(); err != nil {
			fmt.Printf("Failed to close log: %s\n", err)
		}
	}()

	// 初始化数据库
	err = models.InitDB(cfg)
	if err != nil {
		panic(fmt.Sprintf("init db failed, error: %s", err))
	}

	// 启动Web服务
	fmt.Println("Server starting...")
	err = startServer(cfg)
	if err != nil {
		panic(fmt.Sprintf("Server started failed: %s", err))
	}
}

func startServer(cfg *config.AppConfig) error {
	server := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: initEngine(cfg),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go listenToSystemSignals(cancel)

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error(fmt.Sprintf("Failed to shutdown server: %s", err))
		}
	}()
	log.Debug("Server started success")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Debug("Server was shutdown gracefully")
		return nil
	}
	return err
}

func initEngine(cfg *config.AppConfig) *gin.Engine {
	gin.SetMode(func() string {
		if cfg.IsDevEnv() {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "服务器内部错误，请稍后再试！",
		})
	}))
	router.RegisterRoutes(engine)
	return engine
}

func listenToSystemSignals(cancel context.CancelFunc) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-signalChan:
			cancel()
			return
		}
	}
}
