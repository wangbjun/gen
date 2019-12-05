package model

import (
	"gen/config"
	"gen/lib/zlog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/wendal/errors"
	"log"
	"strconv"
)

var DB *gorm.DB

func init() {
	db, err := getDbConnection("default")
	if err != nil {
		log.Println("init mysql pool failed，error：" + err.Error())
	} else {
		DB = db
	}
	log.Println("init mysql pool success")
}

func getDbConnection(name string) (*gorm.DB, error) {
	conf, ok := config.DBConfig[name]
	if !ok {
		return nil, errors.New("database connection [" + name + "] is not existed")
	}
	dsn := conf["username"] + ":" + conf["password"] + "@tcp(" + conf["host"] + ":" + conf["port"] + ")/" +
		conf["database"] + "?charset" + conf["charset"] + "&parseTime=true"
	db, err := gorm.Open(conf["dialect"], dsn)
	if err != nil {
		zlog.Logger.Sugar().Errorf("open database connection failed,error: %s", err.Error())
		return nil, err
	}
	idle, _ := strconv.Atoi(conf["maxIdleConns"])
	open, _ := strconv.Atoi(conf["maxOpenConns"])
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(open)
	if config.Conf.Section("APP").Key("DEBUG").String() == "true" {
		db.LogMode(true)
		db.SetLogger(new(zlog.SqlLog))
	}
	return db, nil
}
