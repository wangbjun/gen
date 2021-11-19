package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"strings"
)

const (
	Dev  = "dev"
	Prod = "prod"
	Test = "test"
)

var Conf *AppConfig

type AppConfig struct {
	File     string
	Env      string
	HttpPort string
	LogMode  string
	LogFile  string
	LogLevel string

	Raw      *ini.File
	DBConfig []DBConfig
}

type DBConfig struct {
	Name        string
	Dsn         string
	MaxIdleConn int
	MaxOpenConn int
}

// Load 加载ini配置文件内容
func Load(file string) (*AppConfig, error) {
	Conf = &AppConfig{
		File:     file,
		Raw:      ini.Empty(),
		Env:      Dev,
		HttpPort: "8080",
	}
	if _, err := os.Stat(Conf.File); os.IsNotExist(err) {
		return nil, fmt.Errorf("cfg file [%s] not existed", Conf.File)
	}
	conf, err := ini.Load(Conf.File)
	if err != nil {
		return nil, fmt.Errorf("load file [%s] failed", Conf.File)
	}
	Conf.Raw = conf
	Conf.loadAppCfg()
	Conf.loadDBCfg()
	return Conf, nil
}

func (cfg *AppConfig) loadAppCfg() {
	section := cfg.Raw.Section("app")

	env := section.Key("env").String()
	if env != "" {
		cfg.Env = env
	}
	httpPort := section.Key("http_port").String()
	if httpPort != "" {
		cfg.HttpPort = httpPort
	}
	logMode := section.Key("log_mode").String()
	if logMode != "" {
		cfg.LogMode = logMode
	}
	logFile := section.Key("log_file").String()
	if logFile != "" {
		cfg.LogFile = logFile
	}
	logLevel := section.Key("log_level").String()
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}
}

func (cfg *AppConfig) loadDBCfg() {
	section := cfg.Raw.Section("db")
	cfg.DBConfig = []DBConfig{
		{
			Name:        "default",
			Dsn:         section.Key("dsn").String(),
			MaxIdleConn: section.Key("max_idle_conn").MustInt(10),
			MaxOpenConn: section.Key("max_open_conn").MustInt(20),
		},
	}
	sections := cfg.Raw.Section("db").ChildSections()
	for _, section := range sections {
		cfg.DBConfig = append(cfg.DBConfig, DBConfig{
			Name:        strings.TrimLeft(section.Name(), "db."),
			Dsn:         section.Key("dsn").String(),
			MaxIdleConn: section.Key("max_idle_conn").MustInt(10),
			MaxOpenConn: section.Key("max_open_conn").MustInt(20),
		})
	}
}

func (cfg *AppConfig) IsDevEnv() bool {
	return cfg.Env == "dev"
}
