package model

import (
	"gen/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	logs "github.com/sirupsen/logrus"
	"github.com/wendal/errors"
	"log"
	"strconv"
)

type Base struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt uint  `json:"created_at"`
	UpdatedAt uint  `json:"updated_at"`
	Status    uint8 `json:"status"`
}

var DB *gorm.DB

func init() {
	db, err := GetDbConnection("default")
	if err != nil {
		log.Println("init mysql pool failed，error：" + err.Error())
	} else {
		DB = db
	}
}

func GetDbConnection(name string) (*gorm.DB, error) {
	conf, ok := config.DBConfig[name]
	if !ok {
		return nil, errors.New("database connection [" + name + "] is not existed")
	}
	dsn := conf["username"] + ":" + conf["password"] + "@tcp(" + conf["host"] + ":" + conf["port"] + ")/" +
		conf["database"] + "?charset" + conf["charset"]
	db, err := gorm.Open(conf["dialect"], dsn)
	if err != nil {
		logs.Errorf("open database connection failed,error: %s", err.Error())
		return nil, err
	}
	idle, _ := strconv.Atoi(conf["maxIdleConns"])
	open, _ := strconv.Atoi(conf["maxOpenConns"])
	db.DB().SetMaxIdleConns(idle)
	db.DB().SetMaxOpenConns(open)
	if config.Conf.String("APP_DEBUG") == "true" {
		db.LogMode(true)
	}
	return db, nil
}
