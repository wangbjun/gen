package model

import (
	"gen/config"
	"gen/zlog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

var dbConnections = make(map[string]*gorm.DB, 0)

func init() {
	for k, v := range config.DBConfig {
		db, err := openConnection(v)
		if err != nil {
			log.Fatalf("init mysql pool [%s] failed，error： %s\n", k, err.Error())
		} else {
			dbConnections[k] = db
			log.Printf("init mysql pool [%s] success\n", k)
		}
	}
}

func DB() *gorm.DB {
	return GetDB("default")
}

func UserDB() *gorm.DB {
	return GetDB("user")
}

func GetDB(name string) *gorm.DB {
	conn, ok := dbConnections[name]
	if !ok {
		return nil
	}
	return conn
}

func openConnection(conf map[string]string) (*gorm.DB, error) {
	db, err := gorm.Open(conf["dialect"], conf["dsn"])
	if err != nil {
		zlog.Logger.Sugar().Errorf("open connection failed,error: %s", err.Error())
		return nil, err
	}
	idle, _ := strconv.Atoi(conf["maxIdleConns"])
	open, _ := strconv.Atoi(conf["maxOpenConns"])
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(open)
	if config.GetAPP("DEBUG").String() == "true" {
		db.LogMode(true)
		db.SetLogger(new(zlog.SqlLog))
	}
	return db, nil
}
