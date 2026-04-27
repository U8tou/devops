package sysauth

import (
	"context"
	"pkg/crypto"
	"pkg/errs"
	sysuser "system/impl/sys_user"
	"system/model"
)

// ChangePassword 修改密码：校验原密码后更新为新密码
func (m *SysAuthImpl) ChangePassword(ctx context.Context, loginId string, oldPassword string, newPassword string) error {
	if oldPassword == newPassword {
		return errs.ERR_NOT_SAME
	}

	userImpl := sysuser.Impl()
	user, err := userImpl.Get(ctx, loginId)
	if err != nil {
		return errs.Sys(err)
	}
	if user == nil {
		return errs.ERR_NOT_ACCOUNT
	}

	if !crypto.CheckPasswordHash(oldPassword, user.Password) {
		return errs.ERR_PWD
	}

	hashedPwd, err := crypto.HashPassword(newPassword)
	if err != nil {
		return errs.Sys(err)
	}

	_, err = m.engine.Context(ctx).Table(&model.SysUser{}).ID(user.Id).Cols("password").Update(&model.SysUser{Password: hashedPwd})
	if err != nil {
		return errs.Sys(err)
	}

	return nil
}
