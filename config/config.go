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
	Raw      *ini.File
	Env      string
	HttpAddr string
	HttpPort string
	LogMode  string
	LogFile  string
	LogLevel string

	DBConfig []dbConfig
}

type dbConfig struct {
	Name        string
	Dialect     string
	Dsn         string
	MaxIdleConn int
	MaxOpenConn int
}

func InitConfig(file string) *AppConfig {
	Conf = &AppConfig{
		File:     file,
		Raw:      ini.Empty(),
		Env:      Dev,
		HttpAddr: "127.0.0.1",
		HttpPort: "8080",
	}
	return Conf
}

// Load 加载ini配置文件内容
func (cfg *AppConfig) Load() error {
	if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
		return fmt.Errorf("cfg file [%s] not existed", cfg.File)
	}
	conf, err := ini.Load(cfg.File)
	if err != nil {
		return fmt.Errorf("load file [%s] failed", cfg.File)
	}
	cfg.Raw = conf
	// load ini config
	cfg.loadAppCfg()
	cfg.loadDBCfg()
	return nil
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
	httpAddr := section.Key("http_addr").String()
	if httpAddr != "" {
		cfg.HttpAddr = httpAddr
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
	cfg.DBConfig = []dbConfig{
		{
			Name:        "default",
			Dialect:     section.Key("dialect").String(),
			Dsn:         section.Key("dsn").String(),
			MaxIdleConn: section.Key("max_idle_conn").MustInt(10),
			MaxOpenConn: section.Key("max_open_conn").MustInt(20),
		},
	}
	sections := cfg.Raw.Section("db").ChildSections()
	for _, section := range sections {
		cfg.DBConfig = append(cfg.DBConfig, dbConfig{
			Name:        strings.TrimLeft(section.Name(), "db."),
			Dialect:     section.Key("dialect").String(),
			Dsn:         section.Key("dsn").String(),
			MaxIdleConn: section.Key("max_idle_conn").MustInt(10),
			MaxOpenConn: section.Key("max_open_conn").MustInt(20),
		})
	}
}

func (cfg *AppConfig) IsDevEnv() bool {
	return cfg.Env == "dev"
}
