package models

import (
	"fmt"
	"gen/config"
	"gen/zlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var (
	db    *gorm.DB
	conns map[string]*gorm.DB
)

func DB(dbName ...string) *gorm.DB {
	if len(dbName) > 0 {
		if conn, ok := conns[dbName[0]]; ok {
			return conn
		}
	}
	return db
}

// Init 初始化数据库连接
func Init(cfg *config.AppConfig) error {
	conns = make(map[string]*gorm.DB)
	for _, v := range cfg.DBConfig {
		conn, err := openConn(v.Dsn, v.MaxIdleConn, v.MaxOpenConn)
		if err != nil {
			return fmt.Errorf("open connection failed, error: %s", err.Error())
		}
		conns[v.Name] = conn
		if v.Name == "default" {
			db = conn
		}
	}
	return nil
}

func openConn(dsn string, idle, open int) (*gorm.DB, error) {
	newLogger := logger.New(Writer{}, logger.Config{SlowThreshold: 500 * time.Millisecond,
		LogLevel: logger.Info, IgnoreRecordNotFoundError: true, Colorful: false})
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

// Writer 记录SQL日志
type Writer struct{}

func (w Writer) Printf(format string, args ...interface{}) {
	zlog.Debug(format, args...)
}
