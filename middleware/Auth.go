package middleware

import (
	"gen/service"
	"github.com/gin-gonic/gin"
	logs "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// 用户鉴权
func Auth() gin.HandlerFunc {
	userService := service.NewUserService()
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userId, err := userService.ParseToken(strings.TrimSpace(strings.Trim(token, "Bearer")))
		if err == nil && userId > 0 {
			logs.Debugf("parse token success, token:%s, userId:%d", token, userId)
			c.Set("userId", userId)
			c.Next()
		} else {
			logs.Errorf("parse token failed, token:%s, error:%s", token, err)
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
		}
	}
}
