package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

var cfg *App

type App struct {
	Env        string
	HttpPort   string
	LogFile    string
	LogConsole bool
	LogLevel   string

	*ini.File
}

func Get() *App {
	return cfg
}

// Init 加载app.ini配置文件
func Init(file string) (*App, error) {
	cfg = &App{
		Env:        Dev,
		HttpPort:   "8080",
		LogConsole: true,
		LogLevel:   "info",
		File:       ini.Empty(),
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, fmt.Errorf("cfg file [%s] not existed", file)
	}
	conf, err := ini.Load(file)
	if err != nil {
		return nil, fmt.Errorf("load file [%s] failed", file)
	}
	cfg.File = conf
	cfg.loadAppCfg()
	return cfg, nil
}

func (cfg *App) IsDevEnv() bool {
	return cfg.Env == Dev
}

func (cfg *App) loadAppCfg() {
	section := cfg.Section("app")
	env := section.Key("env").String()
	if env != "" {
		cfg.Env = env
	}
	httpPort := section.Key("http_port").String()
	if httpPort != "" {
		cfg.HttpPort = httpPort
	}
	logFile := section.Key("log_file").String()
	if logFile != "" {
		cfg.LogFile = logFile
	}
	if section.Key("log_console").String() == "true" {
		cfg.LogConsole = true
	}
	logLevel := section.Key("log_level").String()
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}
}
