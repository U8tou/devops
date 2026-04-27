package sysauth

import (
	"context"
	"pkg/db"
	"system/model"

	"sync"

	"xorm.io/xorm"
)

/**
Notes: SysAuth Impl
Time:  2025-04-29 10:48:59
*/

var (
	sysAuthImpl     *SysAuthImpl
	sysAuthImplOnce sync.Once
)

// SysAuthImpl 登录等业务实现；Token 存储由 pkg/auth 与 auth.use 决定，勿在此拉 Redis 连接。
type SysAuthImpl struct {
	engine *xorm.Engine
}

// Impl 实例化 | 依赖注入
func Impl() ISysAuthImpl {
	sysAuthImplOnce.Do(func() {
		sysAuthImpl = &SysAuthImpl{
			engine: db.GetDb(),
		}
	})
	return sysAuthImpl
}

// ISysAuthImpl 接口
type ISysAuthImpl interface {
	// Login 登入
	Login(ctx context.Context, args *model.LoginReq) (*model.LoginResp, error)
	// RefreshToken 刷新Token（使用RefreshToken换取新的AccessToken）
	RefreshToken(ctx context.Context, device string, refreshToken string) (*model.LoginResp, error)
	// Info 登录信息
	Info(ctx context.Context, loginId string) (*model.InfoResp, error)
	// Logout 登出
	Logout(ctx context.Context, device string, token string) error
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, loginId string, oldPassword string, newPassword string) error
	// UpdateProfile 更新当前用户个人信息（昵称、邮箱、电话、性别、地址）
	UpdateProfile(ctx context.Context, loginId string, nickName, email, phoneArea, phone string, sex int8, address string) error
	// UpdateAvatar 更新当前用户头像地址
	UpdateAvatar(ctx context.Context, loginId string, avatar string) error
	// Register 用户注册（公开接口）
	Register(ctx context.Context, userName string, password string) error
}
