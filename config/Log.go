package config

import (
	"bufio"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func init() {
	logPath := Conf.String("APP_LOG_PATH")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			log.Println("create log dir failed, err: " + err.Error())
		}
	}
	writer, err := rotatelogs.New(
		logPath+"/app.%Y-%m-%d.log",
		rotatelogs.WithMaxAge(time.Hour*24*90),    // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)
	if err != nil {
		log.Printf("config local file system logger error. %+v", errors.WithStack(err))
	}
	setNull()
	lfHook := lfshook.NewHook(writer, &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})
	logrus.AddHook(lfHook)
	// 默认日志记录info级别
	if Conf.String("APP_DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	log.Println("init log config success")
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	logrus.SetOutput(writer)
}
