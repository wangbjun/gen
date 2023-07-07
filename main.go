package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"gen/config"
	"gen/log"
	"gen/middleware"
	"gen/model"
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
	cfg, err := config.Init(configFile)
	if err != nil {
		panic(fmt.Sprintf("load config failed, file: %s, error: %s", configFile, err))
	}

	// 初始化日志
	err = log.Init(cfg)
	if err != nil {
		log.Panic("Init log failed, error: %s", err)
	}
	defer log.Close()

	// 初始化数据库
	err = model.Init(cfg)
	if err != nil {
		log.Panic("Init db failed, error: %s", err)
	}

	// 启动Web服务
	log.Info("Server starting...")
	err = startServer(cfg)
	if err != nil {
		log.Panic("Server started failed: %s", err)
	}
}

func startServer(cfg *config.App) error {
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
			log.Error("Failed to shutdown server: %s", err)
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

func getEngine(cfg *config.App) *gin.Engine {
	gin.SetMode(func() string {
		if cfg.IsDevEnv() {
			return gin.DebugMode
		}
		return gin.ReleaseMode
	}())
	engine := gin.New()
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.WithCtx(c).Error(fmt.Sprintf("server panic: %s", err))
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "服务器内部错误，请稍后再试！",
		})
	}))
	//添加全局的traceId和访问日志
	engine.Use(middleware.AddTrace())
	engine.Use(middleware.AccessLog())
	router.RegisterRoutes(engine)
	return engine
}
