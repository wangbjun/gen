package controller

import (
	"gen/library/zlog"
	"gen/service/userService"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

type userController struct {
	Controller
	userService userService.Service
}

var UserController = &userController{
	userService: userService.New(),
}

// 用户注册
func (uc userController) Register(c *gin.Context) {
	name := c.PostForm("name")
	if !govalidator.StringLength(name, "1", "10") {
		uc.failed(c, ParamError, "名称长度不正确1-10")
		return
	}
	email := c.PostForm("email")
	if !govalidator.IsEmail(email) {
		uc.failed(c, ParamError, "邮箱不正确")
		return
	}
	password := c.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		uc.failed(c, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := uc.userService.Register(name, email, password)
	if err != nil {
		zlog.WithContext(c).Sugar().Errorf("uc register failed, error: %s", err.Error())
		if _, ok := err.(userService.UserError); ok {
			uc.failed(c, ParamError, err.Error())
		} else {
			uc.failed(c, Failed, "注册失败")
		}
	} else {
		zlog.WithContext(c).Sugar().Infof("register uc success, email: %s", email)
		uc.success(c, "ok", map[string]interface{}{"token": token})
	}
	return
}

// 用户登录
func (uc userController) Login(c *gin.Context) {
	email := c.PostForm("email")
	if !govalidator.IsEmail(email) {
		uc.failed(c, ParamError, "邮箱不正确")
		return
	}
	password := c.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		uc.failed(c, ParamError, "密码长度不正确6-16")
		return
	}
	token, err := uc.userService.Login(email, password)
	if err != nil {
		zlog.WithContext(c).Sugar().Errorf("uc register failed, error: %s", err.Error())
		if _, ok := err.(userService.UserError); ok {
			uc.failed(c, Failed, err.Error())
		} else {
			uc.failed(c, Failed, "登录失败")
		}
	} else {
		zlog.WithContext(c).Sugar().Infof("login uc success, email: %s", email)
		uc.success(c, "ok", map[string]interface{}{"token": token})
	}
	return
}

// 用户退出
func (uc userController) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")

	zlog.WithContext(c).Sugar().Debugf("add token into blacklist, token: %s", token)

	uc.success(c, "ok", "")
}
