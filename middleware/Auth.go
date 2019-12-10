package middleware

import (
	"gen/lib/zlog"
	"gen/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 用户鉴权
func Auth() gin.HandlerFunc {
	us := user.New()
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 405,
				"msg":  "未登录",
			})
			return
		}
		userId, err := us.ParseToken(strings.TrimSpace(strings.Trim(token, "Bearer")))
		if err == nil && userId > 0 {
			zlog.WithContext(ctx).Sugar().Debugf("parse token success, userId: %d", userId)
			ctx.Set("userId", userId)
			ctx.Next()
		} else {
			zlog.WithContext(ctx).Sugar().Errorf("parse token failed, error: %s", err)
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 405,
				"msg":  "用户Token无效",
			})
		}
	}
}
