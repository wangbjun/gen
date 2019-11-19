package log

import (
	"gen/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var Sugar *zap.SugaredLogger

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

	Sugar = logger.Sugar()
}

func WithContext(ctx *gin.Context) []interface{} {
	var (
		now          = time.Now().Format("2006-01-02 15:04:05.000")
		serviceStart = ""
		duration     = 0.0
		processId    = os.Getpid()
		request      = ""
		hostAddress  = ""
		clientIp     = ""
		traceId      = ""
		parentId     = ""
		params       = make(map[string][]string)
		size         = 0
	)
	if ctx != nil {
		startTime, _ := ctx.Get("startTime")
		duration = float64(time.Now().Sub(startTime.(time.Time)).Nanoseconds()/1e4) / 100.0 //单位毫秒
		serviceStart = startTime.(time.Time).Format("2006-01-02 15:04:05.000")
		request = ctx.Request.RequestURI
		hostAddress = ctx.Request.Host
		clientIp = ctx.ClientIP()
		traceId = ctx.GetString("traceId")
		parentId = ctx.GetString("parentId")
		params = ctx.Request.PostForm
		size = ctx.Writer.Size()
	}

	return []interface{}{
		"traceId", traceId,
		"serviceStart", serviceStart,
		"serviceEnd", now,
		"processId", processId,
		"request", request,
		"params", params,
		"size", size,
		"hostAddress", hostAddress,
		"clientIp", clientIp,
		"parentId", parentId,
		"duration", duration}
}
