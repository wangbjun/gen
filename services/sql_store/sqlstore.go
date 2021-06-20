package sqlstore

import (
	"fmt"
	"gen/bus"
	"gen/config"
	"gen/log"
	"gen/registry"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"strings"
)

const ServiceName = "SqlStore"
const InitPriority = registry.High

var db *gorm.DB

func init() {
	registry.Register(&registry.Descriptor{
		Name:         ServiceName,
		Instance:     &SQLStore{},
		InitPriority: InitPriority,
	})
}

type SQLStore struct {
	Cfg   *config.Cfg `inject:""`
	Bus   bus.Bus     `inject:""`
	conns map[string]*gorm.DB
	log   *zap.Logger
}

func (ss *SQLStore) Init() error {
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
	return nil
}

func (ss *SQLStore) DB(dbName ...string) *gorm.DB {
	if len(dbName) > 0 {
		if conn, ok := ss.conns[dbName[0]]; ok {
			return conn
		}
	}
	return db
}

func (ss *SQLStore) initDefaultConn() error {
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

func (ss *SQLStore) initChildConns() error {
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

func (ss *SQLStore) openConn(dialect, dsn string, idle, open int) (*gorm.DB, error) {
	conn, err := gorm.Open(dialect, dsn)
	if err != nil {
		return nil, err
	}
	conn.DB().SetMaxIdleConns(idle)
	conn.DB().SetMaxOpenConns(open)
	if ss.Cfg.Env == "dev" {
		conn.LogMode(true)
		conn.SetLogger(new(log.SqlLog))
	}
	return conn, nil
}
