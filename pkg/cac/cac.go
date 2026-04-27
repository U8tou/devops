package cac

import (
	"log"
	"pkg/conf"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	cac     *cache.Cache
	cacOnce sync.Once
)

// 获取实例 | 实例化
func GetCac() *cache.Cache {
	cacOnce.Do(func() {
		cac = new()
	})
	return cac
}

func new() *cache.Cache {
	conf := conf.Cache.Local
	if conf.Expire == 0 {
		conf.Expire = 5
	}
	if conf.Purges == 0 {
		conf.Purges = 10
	}
	ca := cache.New(time.Duration(conf.Expire)*time.Minute, time.Duration(conf.Purges)*time.Minute)
	log.Println("✅ cache ok! ")
	return ca
}
