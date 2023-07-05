package middleware

import (
	"bytes"
	"gen/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"time"
)

func AccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		raw := ctx.Request.URL.RawQuery
		var reqBody []byte
		if ctx.Request.Body != nil {
			reqBody, _ = io.ReadAll(ctx.Request.Body)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
		}
		if len(reqBody) > 1000 {
			reqBody = reqBody[:1000]
		}

		ctx.Next()

		latency := time.Now().Sub(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		bodySize := ctx.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		log.WithCtx(ctx).Info("access_log",
			zap.String("path", path),
			zap.String("method", method),
			zap.String("http_host", ctx.Request.Host),
			zap.String("ua", ctx.Request.UserAgent()),
			zap.String("remote_addr", ctx.Request.RemoteAddr),
			zap.ByteString("request_body", reqBody),
			zap.Int("status_code", statusCode),
			zap.Int("error_code", ctx.GetInt("error_code")),
			zap.String("error_msg", ctx.GetString("error_code")),
			zap.Int("body_size", bodySize),
			zap.String("client_ip", clientIP),
			zap.Duration("latency", latency),
		)
	}
}
