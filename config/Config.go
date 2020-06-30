package config

import (
	"flag"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

var Conf *ini.File

// 默认读取同级目录下的app.ini文件，支持通过 -c 指定配置文件路径
// 为了解决跑测试的时候找不到配置文件的问题，会递归往上找5层目录，如果都找不到，报错
func init() {
	var c string
	var envFile = "./app.ini"
	flag.StringVar(&c, "c", envFile, "custom conf file path")
	flag.Parse()
	if c != "" {
		envFile = c
	} else {
		for i := 0; i < 5; i++ {
			if _, err := os.Stat(envFile); err == nil {
				break
			} else {
				envFile = "../" + envFile
			}
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
