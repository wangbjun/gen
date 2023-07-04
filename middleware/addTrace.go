package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddTrace() gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get("traceId")
		if traceId == "" {
			traceId = uuid.NewString()
		}
		context.Set("traceId", traceId)
		context.Next()
	}
}
