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
		if token == "" {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 401,
				"msg":  "未登录",
			})
		} else {
			userId, err := userService.ParseToken(strings.TrimSpace(strings.Trim(token, "Bearer")))
			if err != nil {
				logs.Errorf("parse token failed，error:%s", err.Error())
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 401,
					"msg":  "未登录",
				})
			}
			if userId <= 0 {
				c.AbortWithStatusJSON(http.StatusOK, gin.H{
					"code": 401,
					"msg":  "未登录",
				})
			} else {
				c.Set("userId", userId)
				c.Next()
			}
		}
	}
}
