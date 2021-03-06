package api

import (
	"context"
	"errors"
	"fmt"
	"gen/config"
	"gen/log"
	"gen/registry"
	"gen/services/article"
	"gen/services/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

const (
	Success      = 200 //正常
	Failed       = 500 //失败
	ParamError   = 400 //参数错误
	NotFound     = 404 //不存在
	UnAuthorized = 401 //未授权
	NotLogin     = 405 //未登录
)

var HttpServer = &HTTPServer{}

func init() {
	registry.Register(&registry.Descriptor{
		Name:         "HTTPServer",
		Instance:     HttpServer,
		InitPriority: registry.High,
	})
}

type HTTPServer struct {
	log     *zap.Logger
	gin     *gin.Engine
	context context.Context

	Cfg            *config.Cfg             `inject:""`
	ArticleService *article.ArticleService `inject:""`
	UserService    *user.UserService       `inject:""`
}

func (hs *HTTPServer) Init() error {
	return nil
}

func (hs *HTTPServer) Run(ctx context.Context) error {
	hs.log = log.New("http_server")

	gin.SetMode(hs.getMode())
	engine := gin.New()
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 500,
			"msg":  "服务器内部错误，请稍后再试！",
		})
	}))

	hs.gin = engine
	hs.context = ctx

	hs.registerRoutes() // 加载路由

	server := &http.Server{Addr: hs.Cfg.HttpAddr + ":" + hs.Cfg.HttpPort, Handler: engine}

	var wg sync.WaitGroup
	wg.Add(1)

	// handle http shutdown on server context done
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			hs.log.Error(fmt.Sprintf("Failed to shutdown server: %s", err))
		}
	}()
	hs.log.Debug("Server was started successfully")
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		hs.log.Debug("Server was shutdown gracefully")
		return nil
	}
	wg.Wait()
	return nil
}

func (hs *HTTPServer) getMode() string {
	debug := hs.Cfg.Env
	if debug == "dev" {
		return gin.DebugMode
	}
	return gin.ReleaseMode
}

func (*HTTPServer) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Gen Web")
}

func (*HTTPServer) Success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": Success,
		"msg":  msg,
		"data": data,
	})
}

func (*HTTPServer) Failed(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": nil,
	})
}
