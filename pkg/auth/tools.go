// Package auth 提供登录尝试、Token、权限、菜单、Session、验证码等能力。
// 使用前须通过 AuthOpt.Build() 完成初始化；获取实例请使用 Get()。
package auth

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ==================== 验证码相关 ====================

// SetCode 设置验证码
func SetCode(ctx context.Context) (string, error) {
	return auth.SetCode(ctx)
}

// CheckCode 校验验证码
func CheckCode(ctx context.Context, id string, code string) bool {
	return auth.CheckCode(ctx, id, code)
}

// ==================== 登录尝试相关 ====================

// TryLogin 尝试登陆，达到阈值后锁定
func TryLogin(ctx context.Context, account string) (*TryLoginResult, error) {
	return auth.TryLogin(ctx, account)
}

// TryClean 清除尝试登陆记录（登录成功后调用）
func TryClean(ctx context.Context, account string) error {
	return auth.TryClean(ctx, account)
}

// IsLoginLocked 检查账号是否被锁定
func IsLoginLocked(ctx context.Context, account string) (bool, time.Duration) {
	return auth.IsLoginLocked(ctx, account)
}

// GetTryCount 获取当前尝试次数
func GetTryCount(ctx context.Context, account string) int64 {
	return auth.GetTryCount(ctx, account)
}

// GetRemainingTries 获取剩余可尝试次数
func GetRemainingTries(ctx context.Context, account string) int64 {
	return auth.GetRemainingTries(ctx, account)
}

// UnlockAccount 手动解锁账号
func UnlockAccount(ctx context.Context, account string) error {
	return auth.UnlockAccount(ctx, account)
}

// LockAccount 手动锁定账号
func LockAccount(ctx context.Context, account string, duration time.Duration) error {
	return auth.LockAccount(ctx, account, duration)
}

// ==================== Token相关 ====================

// GetToken 生成Token（登录时调用）
func GetToken(ctx context.Context, loginId string, device string) TokenInfo {
	return auth.GetToken(ctx, loginId, device)
}

// CheckToken 校验AccessToken，返回loginId（空字符串表示无效）
func CheckToken(ctx context.Context, device string, token string) string {
	return auth.CheckToken(ctx, device, token)
}

// CheckRefreshToken 校验RefreshToken，返回loginId
func CheckRefreshToken(ctx context.Context, device string, token string) string {
	return auth.CheckRefreshToken(ctx, device, token)
}

// RefreshToken 刷新Token（使用RefreshToken换取新的AccessToken）
func RefreshToken(ctx context.Context, device string, refreshToken string) (string, error) {
	return auth.RefreshToken(ctx, device, refreshToken)
}

// RefreshTokenWithRotation 刷新并轮转 Token（旧 RefreshToken 作废，返回新双 Token；推荐用于生产）
func RefreshTokenWithRotation(ctx context.Context, device string, refreshToken string) (TokenInfo, error) {
	return auth.RefreshTokenWithRotation(ctx, device, refreshToken)
}

// OutToken 退出登录（删除AccessToken）
func OutToken(ctx context.Context, device string, token string) error {
	return auth.OutToken(ctx, device, token)
}

// Logout 完全退出登录（同时删除AccessToken和RefreshToken）
func Logout(ctx context.Context, device string, accessToken string, refreshToken string) error {
	return auth.Logout(ctx, device, accessToken, refreshToken)
}

// RenewToken 续期AccessToken（延长过期时间）
func RenewToken(ctx context.Context, device string, token string) error {
	return auth.RenewToken(ctx, device, token)
}

// GetTokenTTL 获取Token剩余有效时间
func GetTokenTTL(ctx context.Context, device string, token string) (time.Duration, error) {
	return auth.GetTokenTTL(ctx, device, token)
}

// ==================== 权限相关 ====================

// SetPermis 设置权限（覆盖）
func SetPermis(ctx context.Context, loginId string, permissions []string) error {
	return auth.SetPermis(ctx, loginId, permissions)
}

// AddPermis 添加权限
func AddPermis(ctx context.Context, loginId string, permissions []string) error {
	return auth.AddPermis(ctx, loginId, permissions)
}

// DelPermis 移除权限
func DelPermis(ctx context.Context, loginId string, permissions []string) error {
	return auth.DelPermis(ctx, loginId, permissions)
}

// GetPermis 获取权限列表
func GetPermis(ctx context.Context, loginId string) ([]string, error) {
	return auth.GetPermis(ctx, loginId)
}

// HasPermis 检查是否拥有所有指定权限
func HasPermis(ctx context.Context, loginId string, permissions ...string) bool {
	return auth.HasPermis(ctx, loginId, permissions...)
}

// HasAnyPermis 检查是否拥有任意一个指定权限
func HasAnyPermis(ctx context.Context, loginId string, permissions []string) bool {
	return auth.HasAnyPermis(ctx, loginId, permissions...)
}

// RootUserType 根用户类型（全量菜单与权限；Auth 中间件跳过权限校验）。与 system/model.UserTypeRoot 一致。
const RootUserType = "00"

// IsRootUser 判断是否为根用户（会话字段 userType=00，由 SysAuth.Login 写入）。旧会话无该字段时需重新登录。
func IsRootUser(ctx context.Context, loginId string) bool {
	v, err := GetSessField(ctx, loginId, "userType")
	if err != nil || v == nil {
		return false
	}
	return strings.TrimSpace(fmt.Sprint(v)) == RootUserType
}

// ==================== 菜单相关 ====================

// SetMenu 设置菜单（覆盖）
func SetMenu(ctx context.Context, loginId string, menus []string) error {
	return auth.SetMenu(ctx, loginId, menus)
}

// AddMenu 添加菜单
func AddMenu(ctx context.Context, loginId string, menus []string) error {
	return auth.AddMenu(ctx, loginId, menus)
}

// DelMenu 移除菜单
func DelMenu(ctx context.Context, loginId string, menus []string) error {
	return auth.DelMenu(ctx, loginId, menus)
}

// GetMenu 获取菜单列表
func GetMenu(ctx context.Context, loginId string) ([]string, error) {
	return auth.GetMenu(ctx, loginId)
}

// HasMenu 检查是否拥有所有指定菜单
func HasMenu(ctx context.Context, loginId string, menus []string) bool {
	return auth.HasMenu(ctx, loginId, menus)
}

// HasAnyMenu 检查是否拥有任意一个指定菜单
func HasAnyMenu(ctx context.Context, loginId string, menus []string) bool {
	return auth.HasAnyMenu(ctx, loginId, menus)
}

// ==================== Session相关 ====================

// SetSess 设置Session（合并）
func SetSess(ctx context.Context, loginId string, mp map[string]any) error {
	return auth.SetSess(ctx, loginId, mp)
}

// GetSess 获取Session
func GetSess(ctx context.Context, loginId string) (map[string]any, error) {
	return auth.GetSess(ctx, loginId)
}

// GetSessField 获取Session单个字段
func GetSessField(ctx context.Context, loginId string, field string) (any, error) {
	return auth.GetSessField(ctx, loginId, field)
}

// DelSessField 删除Session指定字段
func DelSessField(ctx context.Context, loginId string, fields ...string) error {
	return auth.DelSessField(ctx, loginId, fields...)
}

// ClearSess 清空Session
func ClearSess(ctx context.Context, loginId string) error {
	return auth.ClearSess(ctx, loginId)
}

// HasSess 检查Session是否存在
func HasSess(ctx context.Context, loginId string) (bool, error) {
	return auth.HasSess(ctx, loginId)
}

// HasSessField 检查Session字段是否存在
func HasSessField(ctx context.Context, loginId string, field string) (bool, error) {
	return auth.HasSessField(ctx, loginId, field)
}
