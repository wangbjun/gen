package model

import (
	"context"
	"fmt"
	"gen/config"
	"gen/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
)

var (
	connPool = make(map[string]*gorm.DB)
)

// NewOrm 默认返回default数据库连接
func NewOrm(ctx context.Context, dbName ...string) *gorm.DB {
	conn := connPool["default"]
	if len(dbName) > 0 {
		if cn, ok := connPool[dbName[0]]; ok {
			conn = cn
		}
	}
	return conn.WithContext(ctx)
}

// Init 初始化数据库连接
func Init(cfg *config.App) error {
	sections := cfg.Section("db").ChildSections()
	for _, v := range sections {
		var (
			name        = strings.TrimPrefix(v.Name(), "db.")
			dsn         = v.Key("dsn").String()
			maxIdleConn = v.Key("max_idle_conn").MustInt(10)
			maxOpenConn = v.Key("max_open_conn").MustInt(30)
		)
		conn, err := openConn(dsn, maxIdleConn, maxOpenConn)
		if err != nil {
			return fmt.Errorf("open db conn failed, error: %s", err.Error())
		}
		connPool[name] = conn
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
