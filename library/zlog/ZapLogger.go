package zlog

import (
	"gen/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Logger *zap.Logger

func WithContext(ctx *gin.Context) *zap.Logger {
	return Logger.With(getContext(ctx)...)
}

func init() {
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	logFile := config.Conf.Section("APP").Key("LOG_FILE").String()
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    500, // megabytes
		MaxBackups: 0,
		MaxAge:     28, // days
		LocalTime:  true,
	})
	sync := zapcore.AddSync(writer)

	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	core := zapcore.NewTee(zapcore.NewCore(jsonEncoder, sync, infoLevel))

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()

	Logger = logger
}

func getContext(ctx *gin.Context) []zap.Field {
	var (
		now          = time.Now().Format("2006-01-02 15:04:05.000")
		processId    = os.Getpid()
		startTime, _ = ctx.Get("startTime")
		duration     = int(time.Now().Sub(startTime.(time.Time)) / 1e6) //单位毫秒
		serviceStart = startTime.(time.Time).Format("2006-01-02 15:04:05.000")
		request      = ctx.Request.RequestURI
		hostAddress  = ctx.Request.Host
		clientIp     = ctx.ClientIP()
		traceId      = ctx.GetString("traceId")
		parentId     = ctx.GetString("parentId")
		params       = ctx.Request.PostForm
	)
	return []zap.Field{
		zap.String("traceId", traceId),
		zap.String("serviceStart", serviceStart),
		zap.String("serviceEnd", now),
		zap.Int("processId", processId),
		zap.String("request", request),
		zap.String("params", params.Encode()),
		zap.String("hostAddress", hostAddress),
		zap.String("clientIp", clientIp),
		zap.String("parentId", parentId),
		zap.Int("duration", duration)}
}
