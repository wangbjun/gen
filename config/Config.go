package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var Conf *ini.File

func init() {
	envFile := "app.ini"
	// 读取配置文件, 解决跑测试的时候找不到配置文件的问题，最多往上找5层目录
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(envFile); err == nil {
			break
		} else {
			envFile = "../" + envFile
		}
	}
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
