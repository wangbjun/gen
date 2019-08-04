package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	Success      = 200  //正常
	Failed       = 500  //失败
	ParamsError  = 4001 //参数错误
	NotFound     = 4004 //记录不存在
	UnAuthorized = 401  //未授权
	NotLogin     = 405  //未登录
)

type Controller struct{}

var BaseController *Controller

func init() {
	BaseController = &Controller{}
	log.Println("init all controller success")
}

func (*Controller) Index(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Gen Web")
}

func (*Controller) success(ctx *gin.Context, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": Success,
		"msg":  msg,
		"data": data,
	})
}

func (*Controller) failed(ctx *gin.Context, code int, msg string) {
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
	})
}

// 获取当前用户ID
func (*Controller) getUserId(ctx *gin.Context) uint {
	userId, exists := ctx.Get("userId")
	if exists {
		return userId.(uint)
	}
	return 0
}
