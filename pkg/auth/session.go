package auth

import (
	"context"
	"fmt"
	"maps"
	"time"

	"github.com/redis/go-redis/v9"
)

// 設置Session (合并设置，保留未覆盖的字段)
func (m *AuthOpt) SetSess(ctx context.Context, loginId string, mp map[string]any) error {
	if len(mp) == 0 {
		return nil
	}
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		// Redis Hash 仅支持字符串，需显式转换；HSet 需 field-value 对
		vals := make([]any, 0, len(mp)*2)
		for k, v := range mp {
			vals = append(vals, k, sessValToString(v))
		}
		if err := m.rds.HSet(ctx, key, vals...).Err(); err != nil {
			return err
		}
		// 设置过期时间，与 refreshToken 一致
		ttl := time.Duration(m.refTokenTtl) * time.Second
		return m.rds.Expire(ctx, key, ttl).Err()
	}
	// 本地缓存 - 合并已有数据，TTL 与 refreshToken 一致
	ttl := time.Duration(m.refTokenTtl) * time.Second
	existing, has := m.cac.Get(key)
	if has {
		if existMap, ok := existing.(map[string]any); ok {
			maps.Copy(existMap, mp)
			mp = existMap
		}
	}
	m.cac.Set(key, mp, ttl)
	return nil
}

// 獲取Session
func (m *AuthOpt) GetSess(ctx context.Context, loginId string) (map[string]any, error) {
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis HGetAll获取所有字段
		result, err := m.rds.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		// 转换 map[string]string 为 map[string]any
		mp := make(map[string]any, len(result))
		for k, v := range result {
			mp[k] = v
		}
		return mp, nil
	}
	// 本地缓存
	res, has := m.cac.Get(key)
	if !has {
		return map[string]any{}, nil
	}
	if mp, ok := res.(map[string]any); ok {
		return mp, nil
	}
	return map[string]any{}, nil
}

// 獲取Session单个字段
func (m *AuthOpt) GetSessField(ctx context.Context, loginId string, field string) (any, error) {
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		result, err := m.rds.HGet(ctx, key, field).Result()
		if err != nil {
			if err == redis.Nil {
				return nil, nil // key/field 不存在，非异常
			}
			return nil, err
		}
		return result, nil
	}
	// 本地缓存
	res, has := m.cac.Get(key)
	if !has {
		return nil, nil
	}
	if mp, ok := res.(map[string]any); ok {
		return mp[field], nil
	}
	return nil, nil
}

// 刪除Session指定字段
func (m *AuthOpt) DelSessField(ctx context.Context, loginId string, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis HDel删除指定字段
		return m.rds.HDel(ctx, key, fields...).Err()
	}
	// 本地缓存
	res, has := m.cac.Get(key)
	if !has {
		return nil
	}
	mp, ok := res.(map[string]any)
	if !ok {
		return nil
	}
	for _, f := range fields {
		delete(mp, f)
	}
	m.cac.SetDefault(key, mp)
	return nil
}

// 清空Session (删除整个session)
func (m *AuthOpt) ClearSess(ctx context.Context, loginId string) error {
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis Del删除key
		return m.rds.Del(ctx, key).Err()
	}
	// 本地缓存
	m.cac.Delete(key)
	return nil
}

// 检查Session是否存在
func (m *AuthOpt) HasSess(ctx context.Context, loginId string) (bool, error) {
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis Exists检查
		result, err := m.rds.Exists(ctx, key).Result()
		if err != nil {
			return false, err
		}
		return result > 0, nil
	}
	// 本地缓存
	_, has := m.cac.Get(key)
	return has, nil
}

// sessValToString 将 session 值转为字符串（Redis Hash 仅存字符串）
func sessValToString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case int64:
		return fmt.Sprintf("%d", x)
	case int:
		return fmt.Sprintf("%d", x)
	case int32:
		return fmt.Sprintf("%d", x)
	case float64:
		return fmt.Sprintf("%v", x)
	case float32:
		return fmt.Sprintf("%v", x)
	case bool:
		if x {
			return "1"
		}
		return "0"
	default:
		return fmt.Sprint(x)
	}
}

// 检查Session字段是否存在
func (m *AuthOpt) HasSessField(ctx context.Context, loginId string, field string) (bool, error) {
	key := fmt.Sprintf("%s:%s:session", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis HExists检查
		return m.rds.HExists(ctx, key, field).Result()
	}
	// 本地缓存
	res, has := m.cac.Get(key)
	if !has {
		return false, nil
	}
	mp, ok := res.(map[string]any)
	if !ok {
		return false, nil
	}
	_, exists := mp[field]
	return exists, nil
}
