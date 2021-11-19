package log

import (
	"gen/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

// Init 配置日志模块
func Init(cfg *config.AppConfig) {
	var level zapcore.Level
	if level.UnmarshalText([]byte(cfg.LogLevel)) != nil {
		level = zapcore.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "time",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		CallerKey:      "location",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var core zapcore.Core
	switch cfg.LogMode {
	case "console":
		core = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
	case "file":
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogFile,
			MaxSize:    1024, // megabytes
			MaxBackups: 10,
			MaxAge:     30, // days
			LocalTime:  true,
		})
		core = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), level))
	}
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Debug(msg string) {
	Logger.Debug(msg)
}

func Info(msg string) {
	Logger.Info(msg)
}

func Warn(msg string) {
	Logger.Warn(msg)
}

func Error(msg string) {
	Logger.Error(msg)
}

func Panic(msg string) {
	Logger.Panic(msg)
}
