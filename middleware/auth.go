package middleware

import (
	"fmt"
	"gen/log"
	"gen/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// AuthMiddleware 用户鉴权
func AuthMiddleware(user *services.UserService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 405,
				"msg":  "未登录",
			})
			return
		}
		userId, err := user.ParseToken(strings.TrimSpace(strings.Trim(token, "Bearer")))
		if err == nil && userId > 0 {
			log.WithCtx(ctx).Info(fmt.Sprintf("Parse token success, userId: %d", userId))
			ctx.Set("userId", userId)
			ctx.Next()
		} else {
			log.WithCtx(ctx).Error("Parse token failed, error: %s", err)
			ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": 405,
				"msg":  "用户Token无效",
			})
		}
	}
}
