package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

func New(logger string) *zap.Logger {
	return Logger.With(getFields(zap.String("module", logger))...)
}

func init() {
	Logger = zap.NewExample()
}

// Configure 配置日志模块
func Configure(cfg *ini.File) {
	appConfig := cfg.Section("app")
	logMode := appConfig.Key("log_mode").String()
	if logMode == "" {
		logMode = "console"
	}

	logFile := appConfig.Key("log_file").String()
	if logFile == "" {
		logFile = "app.log"
	}

	logLevel := appConfig.Key("log_level").String()
	if logFile == "" {
		logLevel = "info"
	}
	var level zapcore.Level
	if level.UnmarshalText([]byte(logLevel)) != nil {
		level = zapcore.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "name",
		TimeKey:        "time",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var core zapcore.Core
	switch logMode {
	case "console":
		core = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
	case "file":
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500, // megabytes
			MaxBackups: 0,
			MaxAge:     28, // days
			LocalTime:  true,
		})
		core = zapcore.NewTee(zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), level))
	}
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Debug(msg string, fields ...zapcore.Field) {
	Logger.Debug(msg, getFields(fields...)...)
}

func Info(msg string, fields ...zapcore.Field) {
	Logger.Info(msg, getFields(fields...)...)
}

func Warn(msg string, fields ...zapcore.Field) {
	Logger.Warn(msg, getFields(fields...)...)
}

func Error(msg string, fields ...zapcore.Field) {
	Logger.Error(msg, getFields(fields...)...)
}

func Panic(msg string, fields ...zapcore.Field) {
	Logger.Panic(msg, getFields(fields...)...)
}

func getFields(fields ...zapcore.Field) []zapcore.Field {
	var f []zapcore.Field
	if len(fields) > 0 {
		f = append(f, fields...)
	}
	return f
}
