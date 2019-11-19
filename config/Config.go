package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var Conf *ini.File

func init() {
	// 读取配置文件
	envFile := "app.ini"
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Panicf("conf file [%s]  not found!", envFile)
	}
	conf, err := ini.Load(envFile)
	if err != nil {
		log.Panicf("parse conf file [%s] failed, err: %s", envFile, err.Error())
	}

	Conf = conf
	log.Println("init config file success")
}
