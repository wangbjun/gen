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
	File       string
	Env        string
	HttpPort   string
	LogFile    string
	LogConsole bool
	LogLevel   string

	Raw *ini.File
}

type DBConfig struct {
	Name        string
	Dsn         string
	MaxIdleConn int
	MaxOpenConn int
}

// Init 加载app.ini配置文件
func Init(file string) (*App, error) {
	cfg = &App{
		File:     file,
		Raw:      ini.Empty(),
		Env:      Dev,
		HttpPort: "8080",
	}
	if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file [%s] not existed", cfg.File)
	}
	conf, err := ini.Load(cfg.File)
	if err != nil {
		return nil, fmt.Errorf("load file [%s] failed", cfg.File)
	}
	cfg.Raw = conf
	cfg.loadAppCfg()
	return cfg, nil
}

func (cfg *App) IsDevEnv() bool {
	return cfg.Env == Dev
}

func (cfg *App) loadAppCfg() {
	section := cfg.Raw.Section("app")
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
