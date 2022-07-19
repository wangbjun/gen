package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gen/config"
	"gen/models"
	"gen/router"
	"gen/zlog"
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
	zlog.Init(cfg)
	defer zlog.GetLogger().Sync()

	// 初始化数据库
	err = models.Init(cfg)
	if err != nil {
		zlog.Panic("Init db failed, error: %s", err)
	}

	// 启动Web服务
	zlog.Info("Server starting...")
	err = startServer(cfg)
	if err != nil {
		zlog.Panic("Server started failed: %s", err)
	}
}

func startServer(cfg *config.AppConfig) error {
	server := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: getEngine(cfg),
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctxFunc context.CancelFunc) {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-signalChan:
				ctxFunc()
				return
			}
		}
	}(cancel)
	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			zlog.Error("Failed to shutdown server: %s", err)
		}
	}()
	zlog.Debug("Server started success")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		zlog.Debug("Server was shutdown gracefully")
		return nil
	}
	return err
}

func getEngine(cfg *config.AppConfig) *gin.Engine {
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
