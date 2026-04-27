package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/jaevor/go-nanoid"
)

// 设置验证码
func (m *AuthOpt) SetCode(ctx context.Context) (string, error) {
	// 构建验证码对
	decenaryID, err := nanoid.CustomASCII("0123456789", 6)
	if err != nil {
		return "", fmt.Errorf("构建验证码异常: %v", err)
	}

	id := decenaryID()
	code := decenaryID()
	ttl := 60 * time.Second

	codeKey := m.keyPrefix + ":loginCode:" + id

	if m.rds != nil {
		err := m.rds.SetEx(ctx, codeKey, code, ttl).Err() // 60秒过期
		if err != nil {
			return "", fmt.Errorf("构建验证码异常: %v", err)
		}
		return id, nil
	}
	err = m.cac.Add(codeKey, code, ttl)
	if err != nil {
		return "", fmt.Errorf("构建验证码异常: %v", err)
	}
	return id, nil
}

// 校验验证码（验证后立即删除，防止重复使用）
func (m *AuthOpt) CheckCode(ctx context.Context, id string, code string) bool {
	codeKey := m.keyPrefix + ":loginCode:" + id
	if m.rds != nil {
		code_, err := m.rds.GetDel(ctx, codeKey).Result()
		if err != nil {
			return false
		}
		return code_ == code
	}
	code_, ok := m.cac.Get(codeKey)
	if !ok {
		return false
	}
	m.cac.Delete(codeKey)
	s, ok := code_.(string)
	if !ok {
		return false
	}
	return s == code
}
