package middleware

import (
	"gen/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

/**
 * 记录请求日志，加入traceId
 */
func Request() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId, _ := uuid.NewUUID()
		c.Set("traceId", traceId.String())
		c.Set("startTime", time.Now())
		c.Set("parentId", c.GetHeader("X-Ca-Traceid"))
		// before request
		log.Sugar.Infow("请求开始", log.WithContext(c)...)
		c.Next()
		// after request
		log.Sugar.Infow("请求结束", log.WithContext(c)...)
	}
}
