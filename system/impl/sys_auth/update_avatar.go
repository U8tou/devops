package sysauth

import (
	"context"
	"pkg/errs"
	sysuser "system/impl/sys_user"
	"system/model"
)

// UpdateAvatar 更新当前用户头像地址
func (m *SysAuthImpl) UpdateAvatar(ctx context.Context, loginId string, avatar string) error {
	userImpl := sysuser.Impl()
	user, err := userImpl.Get(ctx, loginId)
	if err != nil {
		return errs.Sys(err)
	}
	if user == nil {
		return errs.ERR_NOT_ACCOUNT
	}

	_, err = m.engine.Context(ctx).Table(&model.SysUser{}).ID(user.Id).Cols("avatar").Update(&model.SysUser{Avatar: avatar})
	if err != nil {
		return errs.Sys(err)
	}
	return nil
}
