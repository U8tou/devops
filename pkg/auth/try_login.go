package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// 登录锁定错误（可用 errors.Is 判断）
var ErrLoginLocked = errors.New("账号已锁定，禁止尝试登录")

// LoginLockedError 携带锁定剩余时间的错误，便于上游解析展示
type LoginLockedError struct {
	TTL time.Duration
}

func (e *LoginLockedError) Error() string {
	return fmt.Sprintf("%v，剩余时间：%.0f秒", ErrLoginLocked, e.TTL.Seconds())
}

func (e *LoginLockedError) Unwrap() error { return ErrLoginLocked }

// TryLoginResult 尝试登录结果
type TryLoginResult struct {
	TryCount  int64         // 当前尝试次数
	MaxCount  int64         // 最大允许次数
	Remaining int64         // 剩余可尝试次数
	IsLocked  bool          // 是否已锁定
	LockTTL   time.Duration // 锁定剩余时间
}

const maxAccountLenForTryLogin = 256 // 防止 Redis key 过长或滥用

// 尝试登录 - 记录登录尝试，达到阈值后锁定
// 返回尝试结果和错误（如果已锁定则返回 ErrLoginLocked）
func (m *AuthOpt) TryLogin(ctx context.Context, account string) (*TryLoginResult, error) {
	if len(account) == 0 || len(account) > maxAccountLenForTryLogin {
		return &TryLoginResult{MaxCount: m.RetryCount}, errors.New("账号长度无效")
	}
	key := m.keyPrefix + ":tryLogin:" + account
	lockDuration := time.Duration(m.retryTime) * time.Second

	result := &TryLoginResult{
		MaxCount: m.RetryCount,
	}

	if m.rds != nil {
		return m.tryLoginRedis(ctx, key, lockDuration, result)
	}
	return m.tryLoginCache(ctx, key, lockDuration, result)
}

// Redis 实现（键不存在时 GET 返回 redis.Nil，与业务错误区分）
func (m *AuthOpt) tryLoginRedis(ctx context.Context, key string, lockDuration time.Duration, result *TryLoginResult) (*TryLoginResult, error) {
	tryCount, err := m.rds.Get(ctx, key).Int64()
	if err == nil {
		if tryCount >= m.RetryCount {
			ttl, _ := m.rds.TTL(ctx, key).Result()
			result.TryCount = tryCount
			result.IsLocked = true
			result.LockTTL = ttl
			result.Remaining = 0
			return result, &LoginLockedError{TTL: ttl}
		}
	} else if err != redis.Nil {
		return result, err
	}

	tryCount, err = m.rds.Incr(ctx, key).Result()
	if err != nil {
		return result, err
	}
	if tryCount == 1 {
		_ = m.rds.Expire(ctx, key, lockDuration).Err()
	}

	result.TryCount = tryCount
	result.Remaining = m.RetryCount - tryCount
	if result.Remaining < 0 {
		result.Remaining = 0
	}

	if tryCount >= m.RetryCount {
		_ = m.rds.Expire(ctx, key, lockDuration).Err()
		ttl, _ := m.rds.TTL(ctx, key).Result()
		result.IsLocked = true
		result.LockTTL = ttl
		return result, &LoginLockedError{TTL: ttl}
	}
	return result, nil
}

// 本地缓存实现（接收 ctx 以与 Redis 路径一致，便于后续支持取消）
func (m *AuthOpt) tryLoginCache(ctx context.Context, key string, lockDuration time.Duration, result *TryLoginResult) (*TryLoginResult, error) {
	_ = ctx
	now := time.Now()

	val, expireAt, has := m.cac.GetWithExpiration(key)
	var tryCount int64

	if has {
		if c, ok := val.(int64); ok {
			tryCount = c
		}
		if tryCount >= m.RetryCount {
			ttl := time.Until(expireAt)
			if ttl < 0 {
				ttl = 0
			}
			result.TryCount = tryCount
			result.IsLocked = true
			result.LockTTL = ttl
			result.Remaining = 0
			return result, &LoginLockedError{TTL: ttl}
		}
	}

	// 增加尝试次数
	tryCount++

	// 计算剩余过期时间
	var ttl time.Duration
	if has && !expireAt.IsZero() {
		ttl = expireAt.Sub(now)
		if ttl < 0 {
			ttl = lockDuration
		}
	} else {
		ttl = lockDuration
	}

	// 保存新的尝试次数
	m.cac.Set(key, tryCount, ttl)

	result.TryCount = tryCount
	result.Remaining = m.RetryCount - tryCount
	if result.Remaining < 0 {
		result.Remaining = 0
	}

	if tryCount >= m.RetryCount {
		m.cac.Set(key, tryCount, lockDuration)
		result.IsLocked = true
		result.LockTTL = lockDuration
		return result, &LoginLockedError{TTL: lockDuration}
	}
	return result, nil
}

// 检查账号是否被锁定
func (m *AuthOpt) IsLoginLocked(ctx context.Context, account string) (bool, time.Duration) {
	key := m.keyPrefix + ":tryLogin:" + account

	if m.rds != nil {
		tryCount, err := m.rds.Get(ctx, key).Int64()
		if err != nil {
			return false, 0
		}
		if tryCount >= m.RetryCount {
			ttl, _ := m.rds.TTL(ctx, key).Result()
			return true, ttl
		}
		return false, 0
	}

	// 本地缓存
	val, expireAt, has := m.cac.GetWithExpiration(key)
	if !has {
		return false, 0
	}
	tryCount := val.(int64)
	if tryCount >= m.RetryCount {
		ttl := time.Until(expireAt)
		return true, ttl
	}
	return false, 0
}

// 获取当前尝试次数
func (m *AuthOpt) GetTryCount(ctx context.Context, account string) int64 {
	key := m.keyPrefix + ":tryLogin:" + account

	if m.rds != nil {
		count, err := m.rds.Get(ctx, key).Int64()
		if err != nil {
			return 0
		}
		return count
	}

	// 本地缓存
	val, has := m.cac.Get(key)
	if !has {
		return 0
	}
	return val.(int64)
}

// 获取剩余可尝试次数
func (m *AuthOpt) GetRemainingTries(ctx context.Context, account string) int64 {
	count := m.GetTryCount(ctx, account)
	remaining := m.RetryCount - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// 清除登录尝试记录（登录成功后调用）
func (m *AuthOpt) TryClean(ctx context.Context, account string) error {
	key := m.keyPrefix + ":tryLogin:" + account
	if m.rds != nil {
		return m.rds.Del(ctx, key).Err()
	}
	m.cac.Delete(key)
	return nil
}

// 手动解锁账号
func (m *AuthOpt) UnlockAccount(ctx context.Context, account string) error {
	return m.TryClean(ctx, account)
}

// 手动锁定账号
func (m *AuthOpt) LockAccount(ctx context.Context, account string, duration time.Duration) error {
	key := m.keyPrefix + ":tryLogin:" + account
	if duration <= 0 {
		duration = time.Duration(m.retryTime) * time.Second
	}

	if m.rds != nil {
		// 设置为最大次数，触发锁定
		pipe := m.rds.Pipeline()
		pipe.Set(ctx, key, m.RetryCount, duration)
		_, err := pipe.Exec(ctx)
		return err
	}

	// 本地缓存
	m.cac.Set(key, m.RetryCount, duration)
	return nil
}
