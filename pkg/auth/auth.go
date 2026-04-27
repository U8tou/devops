package auth

import (
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

// 默认配置常量（最佳实践：集中管理默认值）
const (
	DefaultKeyPrefix   = "my_auth"
	DefaultRetryCount  = 3
	DefaultRetryTime   = 60
	DefaultAccTokenTtl = 3 * 60 * 60       // 3 小时
	DefaultRefTokenTtl = 30 * 24 * 60 * 60 // 30 天
)

var (
	auth   *AuthOpt
	authOnce sync.Once
)

type AuthOpt struct {
	rds         *redis.Client
	cac         *cache.Cache
	refTokenTtl int    // 长 token 过期时间(秒)
	accTokenTtl int    // 短 token 过期时间(秒)
	tokenSalt   string // 预留：加密 Salt
	keyPrefix   string // 缓存键前缀
	RetryCount  int64  // 最大可尝试登录次数
	retryTime   int64  // 锁定时长(秒)
}

func (m *AuthOpt) WithStorage(r *redis.Client) *AuthOpt {
	if r != nil {
		m.rds = r
	} else {
		m.cac = cache.New(time.Duration(5)*time.Minute, time.Duration(10)*time.Minute)
	}
	return m
}

func (m *AuthOpt) WithRefTokenTtl(s int) *AuthOpt {
	m.refTokenTtl = s
	return m
}

func (m *AuthOpt) WithAccTokenTtl(s int) *AuthOpt {
	m.accTokenTtl = s
	return m
}

func (m *AuthOpt) WithKeyPrefix(key string) *AuthOpt {
	m.keyPrefix = key
	return m
}

func (m *AuthOpt) WithRetryCount(n int64) *AuthOpt {
	m.RetryCount = n
	return m
}

func (m *AuthOpt) WithRetryTime(s int64) *AuthOpt {
	m.retryTime = s
	return m
}

func (m *AuthOpt) WithTokenSalt(salt string) *AuthOpt {
	m.tokenSalt = salt
	return m
}

// Build 应用默认值并注册为全局单例（仅首次生效，线程安全）
func (m *AuthOpt) Build() {
	authOnce.Do(func() {
		if m.keyPrefix == "" {
			m.keyPrefix = DefaultKeyPrefix
		}
		if m.RetryCount == 0 {
			m.RetryCount = DefaultRetryCount
		}
		if m.retryTime == 0 {
			m.retryTime = DefaultRetryTime
		}
		if m.accTokenTtl == 0 {
			m.accTokenTtl = DefaultAccTokenTtl
		}
		if m.refTokenTtl == 0 {
			m.refTokenTtl = DefaultRefTokenTtl
		}
		auth = m
	})
}

// Get 返回全局 AuthOpt，未初始化时返回 nil（调用方需做 nil 检查）
func Get() *AuthOpt {
	return auth
}
