package controller

import (
	"gen/lib/zlog"
	"gen/service/user"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type userController struct {
	Controller
	userService user.Service
}

var UserController = &userController{
	userService: user.New(),
}

// 用户注册
func (uc userController) Register(ctx *gin.Context) {
	name := ctx.PostForm("name")
	if !govalidator.StringLength(name, "1", "10") {
		uc.failed(ctx, ParamError, "名称长度不正确1-10")
		return
	}
	email := ctx.PostForm("email")
	if !govalidator.IsEmail(email) {
		uc.failed(ctx, ParamError, "邮箱不正确")
		return
	}
	password := ctx.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		uc.failed(ctx, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := uc.userService.Register(name, email, password)
	if err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("uc register failed, error: %s", err.Error())
		if _, ok := err.(user.Error); ok {
			uc.failed(ctx, ParamError, err.Error())
		} else {
			uc.failed(ctx, Failed, "注册失败")
		}
	} else {
		zlog.WithContext(ctx).Sugar().Infof("register uc success, email: %s", email)
		uc.success(ctx, "ok", gin.H{"token": token})
	}
	return
}

// 用户登录
func (uc userController) Login(ctx *gin.Context) {
	email := ctx.PostForm("email")
	if !govalidator.IsEmail(email) {
		uc.failed(ctx, ParamError, "邮箱不正确")
		return
	}
	password := ctx.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		uc.failed(ctx, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := uc.userService.Login(email, password)
	if err != nil {
		zlog.WithContext(ctx).Sugar().Errorf("uc register failed, error: %s", err.Error())
		if _, ok := err.(user.Error); ok {
			uc.failed(ctx, Failed, err.Error())
		} else {
			uc.failed(ctx, Failed, "登录失败")
		}
	} else {
		zlog.WithContext(ctx).Sugar().Infof("login uc success, email: %s", email)
		uc.success(ctx, "ok", gin.H{"token": token})
	}
	return
}

// 用户退出
func (uc userController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	zlog.WithContext(ctx).Sugar().Debugf("add token into blacklist, token: %s", token)

	uc.success(ctx, "ok", "")
}
