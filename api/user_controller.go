package api

import (
	"fmt"
	"gen/log"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type userController struct {
	*HTTPServer
}

var UserController = &userController{
	httpServer,
}

// Register 用户注册
func (r userController) Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	if !govalidator.StringLength(name, "1", "10") {
		r.Failed(ctx, ParamError, "名称长度不正确1-10")
		return
	}
	email := ctx.PostForm("email")
	if !govalidator.IsEmail(email) {
		r.Failed(ctx, ParamError, "邮箱不正确")
		return
	}
	password := ctx.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		r.Failed(ctx, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := r.HTTPServer.UserService.Register(name, email, password)
	if err != nil {
		log.Error(fmt.Sprintf("r register Failed, error: %s", err.Error()))
		r.Failed(ctx, Failed, "注册失败")
	} else {
		log.Info(fmt.Sprintf("register r Success, email: %s", email))
		r.Success(ctx, "ok", gin.H{"token": token})
	}
	return
}

// Login 用户登录
func (r userController) Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	if !govalidator.IsEmail(email) {
		r.Failed(ctx, ParamError, "邮箱不正确")
		return
	}
	password := ctx.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		r.Failed(ctx, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := r.HTTPServer.UserService.Login(email, password)
	if err != nil {
		log.Error(fmt.Sprintf("r register Failed, error: %s", err.Error()))
		r.Failed(ctx, Failed, "登录失败")
	} else {
		log.Info(fmt.Sprintf("login r Success, email: %s", email))
		r.Success(ctx, "ok", gin.H{"token": token})
	}
	return
}

// Logout 用户退出
func (r userController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	log.Debug(fmt.Sprintf("add token into blacklist, token: %s", token))

	r.Success(ctx, "ok", "")
}
