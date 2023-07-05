package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddTrace() gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get("trace_id")
		if traceId == "" {
			traceId = uuid.NewString()
		}
		context.Set("trace_id", traceId)
		context.Next()
	}
}
