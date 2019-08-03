package model

import (
	. "gen/config"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type Base struct {
	ID        uint  `json:"id" gorm:"primary_key"`
	CreatedAt uint  `json:"created_at"`
	UpdatedAt uint  `json:"updated_at"`
	Status    uint8 `json:"status"`
}

var DB *gorm.DB

var Redis *redis.Client

func init() {
	initDB()
	initRedis()
}

// 初始化DB
func initDB() {
	var (
		dialect      = Conf.DefaultString("DB_Dialect", "mysql")
		host         = Conf.DefaultString("DB_HOST", "127.0.0.1")
		port         = Conf.DefaultString("DB_PORT", "3306")
		user         = Conf.DefaultString("DB_USERNAME", "user")
		pass         = Conf.DefaultString("DB_PASSWORD", "pass")
		database     = Conf.String("DB_DATABASE")
		charset      = Conf.DefaultString("DB_CHARSET", "utf8mb4")
		maxIdleConns = Conf.DefaultInt("DB_MAX_IDLE_CONN", 5)
		maxOpenConns = Conf.DefaultInt("DB_MAX_OPEN_CONN", 50)
	)
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + database + "?charset" + charset
	db, err := gorm.Open(dialect, dsn)
	if err != nil {
		log.Fatalf("init db connection failed, error: " + err.Error())
	}
	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	if Conf.String("APP_DEBUG") == "true" {
		db.LogMode(true)
	}
	DB = db
	log.Println("init db connection success")
}

// 初始化redis
func initRedis() {
	var (
		host   = Conf.DefaultString("REDIS_HOST", "127.0.0.1")
		port   = Conf.DefaultString("REDIS_PORT", "3306")
		pass   = Conf.String("REDIS_PASS")
		minIde = Conf.DefaultInt("REDIS_MIN_IDLE", 5)
	)
	client := redis.NewClient(
		&redis.Options{
			Addr:         host + ":" + port,
			Password:     pass,
			MaxRetries:   3,
			MinIdleConns: minIde,
		},
	)
	Redis = client
	log.Println("init redis pool success")
}
