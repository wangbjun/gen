package middleware

import (
	"gen/library/zlog"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

/**
 * 记录请求日志，加入traceId
 */
func Request() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId, _ := uuid.NewUUID()
		ctx.Set("traceId", traceId.String())
		ctx.Set("startTime", time.Now())
		ctx.Set("parentId", ctx.GetHeader("X-Ca-Traceid"))
		// before request
		zlog.WithContext(ctx).Info("请求开始")
		ctx.Next()
		// after request
		zlog.WithContext(ctx).Info("请求结束")
	}
}
