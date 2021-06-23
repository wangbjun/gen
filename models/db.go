package models

import (
	"fmt"
	"gen/config"
	"gen/log"
	"gen/registry"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strings"
	"time"
)

var (
	db       *gorm.DB
	sqlStore *SQLService
)

func init() {
	registry.Register(&registry.Descriptor{
		Name:         "SqlService",
		Instance:     &SQLService{},
		InitPriority: registry.High,
	})
}

type SQLService struct {
	Cfg *config.Cfg `inject:""`

	conns map[string]*gorm.DB
	log   *zap.Logger
}

func DB(dbName ...string) *gorm.DB {
	if len(dbName) > 0 {
		if conn, ok := sqlStore.conns[dbName[0]]; ok {
			return conn
		}
	}
	return db
}

func (ss *SQLService) Init() error {
	ss.log = log.Logger
	ss.conns = make(map[string]*gorm.DB)
	if err := ss.initDefaultConn(); err != nil {
		ss.log.Error(fmt.Sprintf("init default db conn failed: %s", err.Error()))
		return err
	}

	if err := ss.initChildConns(); err != nil {
		ss.log.Error(fmt.Sprintf("init child db conn failed: %s", err.Error()))
		return err
	}
	sqlStore = ss
	return nil
}

func (ss *SQLService) initDefaultConn() error {
	section := ss.Cfg.Raw.Section("db")
	var (
		dialect         = section.Key("dialect").String()
		dsn             = section.Key("dsn").String()
		maxIdleConns, _ = section.Key("max_idle_conn").Int()
		maxOpenConns, _ = section.Key("max_open_conn").Int()
	)
	conn, err := ss.openConn(dialect, dsn, maxIdleConns, maxOpenConns)
	if err != nil {
		return fmt.Errorf("open connection failed, error: %s", err.Error())
	}
	ss.conns["default"] = conn
	db = conn
	return nil
}

func (ss *SQLService) initChildConns() error {
	sections := ss.Cfg.Raw.Section("db").ChildSections()
	for _, section := range sections {
		var (
			dialect         = section.Key("dialect").String()
			dsn             = section.Key("dsn").String()
			maxIdleConns, _ = section.Key("max_idle_conn").Int()
			maxOpenConns, _ = section.Key("max_open_conn").Int()
		)
		conn, err := ss.openConn(dialect, dsn, maxIdleConns, maxOpenConns)
		if err != nil {
			return fmt.Errorf("open connection failed, error: %s", err.Error())
		}
		ss.conns[strings.TrimLeft(section.Name(), "db.")] = conn
	}
	return nil
}

func (ss *SQLService) openConn(dialect, dsn string, idle, open int) (*gorm.DB, error) {
	newLogger := logger.New(Writer{}, logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false})
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
	log.Info(fmt.Sprintf(format, args...))
}
