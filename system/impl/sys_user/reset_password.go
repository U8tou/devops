package sysuser

import (
	"context"
	"pkg/crypto"
	"system/model"
)

// ResetPassword 重置密码（管理员直接设置新密码）
func (m *SysUserImpl) ResetPassword(ctx context.Context, userId int64, newPassword string) (int64, error) {
	hashedPwd, err := crypto.HashPassword(newPassword)
	if err != nil {
		return 0, err
	}
	return m.engine.Context(ctx).
		Table(&model.SysUser{}).
		ID(userId).
		Cols("password").
		Update(&model.SysUser{Password: hashedPwd})
}

