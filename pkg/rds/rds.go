package rds

import (
	"context"
	"fmt"
	"log"
	"pkg/conf"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	rds     *redis.Client
	rdsOnce sync.Once
)

// 获取实例 | 实例化
func GetRds() *redis.Client {
	rdsOnce.Do(func() {
		rds = new()
	})
	return rds
}

func new() *redis.Client {
	conf := conf.Cache.Redis
	if conf.Port == 0 {
		conf.Port = 6379
	}
	if conf.Addr == "" {
		conf.Addr = "127.0.0.1"
	}
	r := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Addr, conf.Port),
		Password: conf.Pwd,
		DB:       conf.Db,
	})
	if _, err := r.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}
	fmt.Println("✅ redis ok! ")
	return r
}
