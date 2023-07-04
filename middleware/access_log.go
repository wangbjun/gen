package middleware

import (
	"gen/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery

		ctx.Next()

		latency := time.Now().Sub(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		requestBody := ctx.Request.Body
		statusCode := ctx.Writer.Status()
		errorMessage := ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		bodySize := ctx.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		log.WithCtx(ctx).Info("access_log",
			zap.String("path", path),
			zap.String("method", method),
			zap.Any("request_body", requestBody),
			zap.Int("status_code", statusCode),
			zap.String("error", errorMessage),
			zap.Int("body_size", bodySize),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}
