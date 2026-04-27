package auth

import (
	"context"
	"fmt"
	"time"
)

// 設置權限 (覆盖设置，先删除旧权限再添加新权限；Redis 下设置 TTL 与 Session 一致)
func (m *AuthOpt) SetPermis(ctx context.Context, loginId string, ps []string) error {
	key := fmt.Sprintf("%s:%s:permis", m.keyPrefix, loginId)
	if m.rds != nil {
		pipe := m.rds.Pipeline()
		pipe.Del(ctx, key)
		if len(ps) > 0 {
			pipe.SAdd(ctx, key, convertToAny(ps)...)
			ttl := time.Duration(m.refTokenTtl) * time.Second
			if ttl > 0 {
				pipe.Expire(ctx, key, ttl)
			}
		}
		_, err := pipe.Exec(ctx)
		return err
	}
	m.cac.SetDefault(key, ps)
	return nil
}

// 添加權限
func (m *AuthOpt) AddPermis(ctx context.Context, loginId string, ps []string) error {
	key := fmt.Sprintf("%s:%s:permis", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis
		return m.rds.SAdd(ctx, key, convertToAny(ps)...).Err()
	}
	// 本地缓存
	existing, has := m.cac.Get(key)
	if has {
		existList := existing.([]string)
		ps = append(existList, ps...)
	}
	m.cac.SetDefault(key, ps)
	return nil
}

// 移除權限
func (m *AuthOpt) DelPermis(ctx context.Context, loginId string, ps []string) error {
	key := fmt.Sprintf("%s:%s:permis", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis
		return m.rds.SRem(ctx, key, convertToAny(ps)...).Err()
	}
	// 本地缓存
	existing, has := m.cac.Get(key)
	if !has {
		return nil
	}
	existList := existing.([]string)
	newList := make([]string, 0, len(existList))
	removeSet := make(map[string]struct{}, len(ps))
	for _, p := range ps {
		removeSet[p] = struct{}{}
	}
	for _, p := range existList {
		if _, ok := removeSet[p]; !ok {
			newList = append(newList, p)
		}
	}
	m.cac.SetDefault(key, newList)
	return nil
}

// 獲取權限列表
func (m *AuthOpt) GetPermis(ctx context.Context, loginId string) ([]string, error) {
	key := fmt.Sprintf("%s:%s:permis", m.keyPrefix, loginId)
	if m.rds != nil {
		// 使用redis
		return m.rds.SMembers(ctx, key).Result()
	}
	// 本地缓存
	existing, has := m.cac.Get(key)
	if !has {
		return []string{}, nil
	}
	return existing.([]string), nil
}

// 檢查權限 (检查用户是否拥有所有指定的权限；Redis 下使用 SMIsMember 一次请求)
func (m *AuthOpt) HasPermis(ctx context.Context, loginId string, permissions ...string) bool {
	if len(permissions) == 0 {
		return true
	}
	key := fmt.Sprintf("%s:%s:permis", m.keyPrefix, loginId)
	if m.rds != nil {
		vals, err := m.rds.SMIsMember(ctx, key, convertToAny(permissions)...).Result()
		if err != nil || len(vals) != len(permissions) {
			return false
		}
		for _, v := range vals {
			if !v {
				return false
			}
		}
		return true
	}
	existing, has := m.cac.Get(key)
	if !has {
		return false
	}
	existList, _ := existing.([]string)
	permisSet := make(map[string]struct{}, len(existList))
	for _, p := range existList {
		permisSet[p] = struct{}{}
	}
	for _, p := range permissions {
		if _, ok := permisSet[p]; !ok {
			return false
		}
	}
	return true
}

// 檢查權限 (检查用户是否拥有任意一个指定的权限；Redis 下使用 SMIsMember 一次请求)
func (m *AuthOpt) HasAnyPermis(ctx context.Context, loginId string, permissions ...string) bool {
	if len(permissions) == 0 {
		return true
	}
	key := fmt.Sprintf("%s:%s:permis", m.keyPrefix, loginId)
	if m.rds != nil {
		vals, err := m.rds.SMIsMember(ctx, key, convertToAny(permissions)...).Result()
		if err != nil {
			return false
		}
		for _, v := range vals {
			if v {
				return true
			}
		}
		return false
	}
	existing, has := m.cac.Get(key)
	if !has {
		return false
	}
	existList, _ := existing.([]string)
	permisSet := make(map[string]struct{}, len(existList))
	for _, p := range existList {
		permisSet[p] = struct{}{}
	}
	for _, p := range permissions {
		if _, ok := permisSet[p]; ok {
			return true
		}
	}
	return false
}

// 辅助函数：将 []string 转换为 []any
func convertToAny(ss []string) []any {
	result := make([]any, len(ss))
	for i, s := range ss {
		result[i] = s
	}
	return result
}
