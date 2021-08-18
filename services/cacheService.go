package services

import (
	"gen/config"
	"gen/registry"
	"github.com/go-redis/redis/v8"
)

type CacheService struct {
	Cfg   *config.Cfg `inject:""`
	Redis *redis.Client
}

func init() {
	registry.RegisterService(&CacheService{})
}

func (r *CacheService) Init() error {
	cfg := r.Cfg.Raw.Section("redis")
	var (
		host       = cfg.Key("host").String()
		port       = cfg.Key("port").String()
		pass       = cfg.Key("pass").String()
		minIdle, _ = cfg.Key("min_idle").Int()
	)
	r.Redis = redis.NewClient(&redis.Options{
		Addr:         host + ":" + port,
		Password:     pass,
		DB:           0,
		MinIdleConns: minIdle,
	})
	return nil
}
