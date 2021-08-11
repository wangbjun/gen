package controllers

import (
	"fmt"
	"gen/log"
	"gen/models"
	"gen/utils/trans"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type userController struct {
	*HTTPServer
}

// Register 用户注册
func (r userController) Register(ctx *gin.Context) {
	var form models.UserRegisterCommand
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			r.Failed(ctx, ParamError, trans.Translate(e))
		} else {
			r.Failed(ctx, Failed, "请求错误")
		}
		return
	}
	token, err := r.HTTPServer.UserService.Register(&form)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
		r.Success(ctx, "ok", gin.H{"token": token})
	}
	return
}

// Login 用户登录
func (r userController) Login(ctx *gin.Context) {
	var form models.UserLoginCommand
	err := ctx.ShouldBindJSON(&form)
	if err != nil {
		if e, ok := err.(validator.ValidationErrors); ok {
			r.Failed(ctx, ParamError, trans.Translate(e))
		} else {
			r.Failed(ctx, Failed, "请求错误")
		}
		return
	}
	token, err := r.HTTPServer.UserService.Login(&form)
	if err != nil {
		r.Failed(ctx, Failed, err.Error())
	} else {
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
