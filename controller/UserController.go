package controller

import (
	"gen/service"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
)

type userController struct {
	Controller
	userService service.UserService
}

var UserController = &userController{
	userService: service.NewUserService(),
}

// 用户注册
func (u userController) Register(c *gin.Context) {
	name := c.PostForm("name")
	if !govalidator.StringLength(name, "1", "10") {
		u.failed(c, ParamsError, "名称长度不正确1-10")
		return
	}
	email := c.PostForm("email")
	if !govalidator.IsEmail(email) {
		u.failed(c, ParamsError, "邮箱不正确")
		return
	}
	password := c.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		u.failed(c, ParamsError, "密码长度不正确6-16")
		return
	}
	token, err := u.userService.Register(name, email, password)
	if err != nil {
		if err == service.UserExisted {
			logs.Errorf("user register failed, error: %s", err.Error())
			u.failed(c, ParamsError, err.Error())
		} else {
			logs.Errorf("user register failed, error: %s", err.Error())
			u.failed(c, Failed, "注册失败")
		}
	} else {
		logs.Infof("register user success, email:%s", email)
		u.success(c, "ok", map[string]interface{}{"token": token})
	}
	return
}

// 用户登录
func (u userController) Login(c *gin.Context) {
	email := c.PostForm("email")
	if !govalidator.IsEmail(email) {
		u.failed(c, ParamsError, "邮箱不正确")
		return
	}
	password := c.PostForm("password")
	if !govalidator.StringLength(password, "6", "16") {
		u.failed(c, ParamsError, "密码长度不正确6-16")
		return
	}
	token, err := u.userService.Login(email, password)
	if err != nil {
		if err == service.UserNotExisted || err == service.PasswordWrong {
			u.failed(c, Failed, err.Error())
		} else {
			u.failed(c, Failed, "登录失败")
		}
	} else {
		logs.Infof("login user success, email:%s", email)
		u.success(c, "ok", map[string]interface{}{"token": token})
	}
	return
}

// 用户退出
func (u userController) Logout(c *gin.Context) {

}
