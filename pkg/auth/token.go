package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jaevor/go-nanoid"
)

// Token 相关错误
var (
	ErrTokenInvalid  = errors.New("token无效或已过期")
	ErrTokenGenerate = errors.New("token生成失败")
)

// TokenInfo Token信息
type TokenInfo struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	ExpiresIn    int    `json:"expires_in"`    // 访问令牌过期时间(秒)
}

// 生成Token（登录时调用）
func (m *AuthOpt) GetToken(ctx context.Context, loginId string, device string) TokenInfo {
	// 生成 AccessToken
	accessKey := m.keyPrefix + ":accessToken:" + device + ":"
	accessTTL := time.Duration(m.accTokenTtl) * time.Second
	accessToken := m.createToken(ctx, loginId, accessTTL, accessKey)

	// 生成 RefreshToken（使用更长的过期时间）
	refreshKey := m.keyPrefix + ":refreshToken:" + device + ":"
	refreshTTL := time.Duration(m.refTokenTtl) * time.Second
	refreshToken := m.createToken(ctx, loginId, refreshTTL, refreshKey)

	return TokenInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    m.accTokenTtl,
	}
}

// 校验AccessToken，返回loginId
func (m *AuthOpt) CheckToken(ctx context.Context, device string, token string) string {
	key := m.keyPrefix + ":accessToken:" + device + ":" + token
	return m.getTokenLoginId(ctx, key)
}

// 校验RefreshToken，返回loginId
func (m *AuthOpt) CheckRefreshToken(ctx context.Context, device string, token string) string {
	key := m.keyPrefix + ":refreshToken:" + device + ":" + token
	return m.getTokenLoginId(ctx, key)
}

// RefreshToken 使用 RefreshToken 换取新的 AccessToken（不轮转，原 RefreshToken 仍有效）
func (m *AuthOpt) RefreshToken(ctx context.Context, device string, refreshToken string) (string, error) {
	loginId := m.CheckRefreshToken(ctx, device, refreshToken)
	if loginId == "" {
		return "", ErrTokenInvalid
	}
	accessKey := m.keyPrefix + ":accessToken:" + device + ":"
	accessTTL := time.Duration(m.accTokenTtl) * time.Second
	newAccessToken := m.createToken(ctx, loginId, accessTTL, accessKey)
	if newAccessToken == "" {
		return "", ErrTokenGenerate
	}
	return newAccessToken, nil
}

// RefreshTokenWithRotation 刷新并轮转 Token（最佳实践：旧 RefreshToken 作废，返回新双 Token）
func (m *AuthOpt) RefreshTokenWithRotation(ctx context.Context, device string, refreshToken string) (TokenInfo, error) {
	loginId := m.CheckRefreshToken(ctx, device, refreshToken)
	if loginId == "" {
		return TokenInfo{}, ErrTokenInvalid
	}
	refreshKey := m.keyPrefix + ":refreshToken:" + device + ":" + refreshToken
	_ = m.deleteKey(ctx, refreshKey)

	info := m.GetToken(ctx, loginId, device)
	return info, nil
}

// 退出登录（删除AccessToken）
func (m *AuthOpt) OutToken(ctx context.Context, device string, token string) error {
	key := m.keyPrefix + ":accessToken:" + device + ":" + token
	return m.deleteKey(ctx, key)
}

// 完全退出登录（同时删除AccessToken和RefreshToken）
func (m *AuthOpt) Logout(ctx context.Context, device string, accessToken string, refreshToken string) error {
	accessKey := m.keyPrefix + ":accessToken:" + device + ":" + accessToken
	refreshKey := m.keyPrefix + ":refreshToken:" + device + ":" + refreshToken

	if m.rds != nil {
		// 使用 Pipeline 批量删除
		pipe := m.rds.Pipeline()
		pipe.Del(ctx, accessKey)
		pipe.Del(ctx, refreshKey)
		_, err := pipe.Exec(ctx)
		return err
	}

	// 本地缓存
	m.cac.Delete(accessKey)
	m.cac.Delete(refreshKey)
	return nil
}

// 续期AccessToken（延长过期时间）
func (m *AuthOpt) RenewToken(ctx context.Context, device string, token string) error {
	key := m.keyPrefix + ":accessToken:" + device + ":" + token
	ttl := time.Duration(m.accTokenTtl) * time.Second

	if m.rds != nil {
		// 检查 token 是否存在
		exists, err := m.rds.Exists(ctx, key).Result()
		if err != nil || exists == 0 {
			return ErrTokenInvalid
		}
		return m.rds.Expire(ctx, key, ttl).Err()
	}

	// 本地缓存
	val, has := m.cac.Get(key)
	if !has {
		return ErrTokenInvalid
	}
	m.cac.Set(key, val, ttl)
	return nil
}

// 获取Token剩余有效时间
func (m *AuthOpt) GetTokenTTL(ctx context.Context, device string, token string) (time.Duration, error) {
	key := m.keyPrefix + ":accessToken:" + device + ":" + token

	if m.rds != nil {
		ttl, err := m.rds.TTL(ctx, key).Result()
		if err != nil {
			return 0, err
		}
		if ttl < 0 {
			return 0, ErrTokenInvalid
		}
		return ttl, nil
	}

	// 本地缓存
	_, expireAt, has := m.cac.GetWithExpiration(key)
	if !has {
		return 0, ErrTokenInvalid
	}
	ttl := time.Until(expireAt)
	if ttl < 0 {
		return 0, ErrTokenInvalid
	}
	return ttl, nil
}

// ==================== 内部辅助方法 ====================

// 创建令牌
func (m *AuthOpt) createToken(ctx context.Context, loginId string, ttl time.Duration, preKey string) string {
	gen, err := nanoid.Canonic()
	if err != nil {
		return ""
	}
	token := gen()
	key := preKey + token

	if m.rds != nil {
		if err := m.rds.SetEx(ctx, key, loginId, ttl).Err(); err != nil {
			return ""
		}
		return token
	}

	// 本地缓存
	if err := m.cac.Add(key, loginId, ttl); err != nil {
		// Add 失败可能是 key 已存在，尝试 Set
		m.cac.Set(key, loginId, ttl)
	}
	return token
}

// 获取token对应的loginId
func (m *AuthOpt) getTokenLoginId(ctx context.Context, key string) string {
	if m.rds != nil {
		loginId, err := m.rds.Get(ctx, key).Result()
		if err != nil {
			return ""
		}
		return loginId
	}

	// 本地缓存
	loginId, has := m.cac.Get(key)
	if !has {
		return ""
	}
	if s, ok := loginId.(string); ok {
		return s
	}
	return ""
}

// 删除key
func (m *AuthOpt) deleteKey(ctx context.Context, key string) error {
	if m.rds != nil {
		return m.rds.Del(ctx, key).Err()
	}
	m.cac.Delete(key)
	return nil
}
