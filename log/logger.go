package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	"gen/config"
)

var zapLogger *zap.Logger

type Logger struct {
	context.Context
	*zap.Logger
}

// WithCtx 带请求上下文的Logger，可以记录一些额外信息，比如traceId
func WithCtx(ctx context.Context) *Logger {
	return &Logger{ctx, zapLogger}
}

func Close() {
	zapLogger.Sync()
}

// Init 初始化配置日志模块
func Init(cfg *config.App) error {
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
	var cores []zapcore.Core
	if cfg.LogFile != "" {
		fileWriter, err := os.OpenFile(cfg.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriter), level))
	}
	if cfg.LogConsole {
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level))
	}
	zapLogger = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	return nil
}

func Debug(format string, args ...interface{}) {
	zapLogger.Info(fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	zapLogger.Info(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	zapLogger.Warn(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	zapLogger.Error(fmt.Sprintf(format, args...))
}

func Panic(format string, args ...interface{}) {
	zapLogger.Panic(fmt.Sprintf(format, args...))
}

func (r Logger) Debug(format string, args ...interface{}) {
	traceId, _ := r.Value("trace_id").(string)
	r.Logger.With(zap.String("trace_id", traceId)).Debug(fmt.Sprintf(format, args...))
}

func (r Logger) Info(msg string, fields ...zap.Field) {
	traceId, _ := r.Value("trace_id").(string)
	r.Logger.With(append(fields, zap.String("trace_id", traceId))...).Info(msg)
}

func (r Logger) Warn(format string, args ...interface{}) {
	traceId, _ := r.Value("trace_id").(string)
	r.Logger.With(zap.String("trace_id", traceId)).Warn(fmt.Sprintf(format, args...))
}

func (r Logger) Error(format string, args ...interface{}) {
	traceId, _ := r.Value("trace_id").(string)
	r.Logger.With(zap.String("trace_id", traceId)).Error(fmt.Sprintf(format, args...))
}

func (r Logger) Panic(format string, args ...interface{}) {
	traceId, _ := r.Value("trace_id").(string)
	r.Logger.With(zap.String("trace_id", traceId)).Panic(fmt.Sprintf(format, args...))
}
