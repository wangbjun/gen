package models

import (
	"fmt"
	"gen/config"
	"gen/log"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

var (
	db    *gorm.DB
	conns = make(map[string]*gorm.DB)
)

func NewDB(dbName ...string) *gorm.DB {
	if len(dbName) > 0 {
		if conn, ok := conns[dbName[0]]; ok {
			return conn
		}
	}
	return db
}

// Init 初始化数据库连接
func Init(cfg *config.App) error {
	sections := []*ini.Section{cfg.Raw.Section("db")}
	sections = append(sections, cfg.Raw.Section("db").ChildSections()...)
	for _, v := range sections {
		var (
			name        = strings.TrimLeft(v.Name(), "db.")
			dsn         = v.Key("dsn").String()
			maxIdleConn = v.Key("max_idle_conn").MustInt(10)
			maxOpenConn = v.Key("max_open_conn").MustInt(30)
		)
		conn, err := openConn(dsn, maxIdleConn, maxOpenConn)
		if err != nil {
			return fmt.Errorf("open db conn failed, error: %s", err.Error())
		}
		if name == "" {
			db = conn
		} else {
			conns[name] = conn
		}
	}
	return nil
}

func openConn(dsn string, idle, open int) (*gorm.DB, error) {
	newLogger := log.NewGormLogger(logger.Config{
		LogLevel: logger.Info,
	})
	openDB, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}
	db, err := openDB.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(idle)
	db.SetMaxOpenConns(open)
	return openDB, nil
}
