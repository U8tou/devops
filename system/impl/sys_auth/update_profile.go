package sysauth

import (
	"context"
	"pkg/auth"
	"pkg/errs"
	sysuser "system/impl/sys_user"
	"system/model"
)

// UpdateProfile 更新当前用户个人信息（仅允许昵称、邮箱、电话、性别、地址）
func (m *SysAuthImpl) UpdateProfile(ctx context.Context, loginId string, nickName, email, phoneArea, phone string, sex int8, address string) error {
	userImpl := sysuser.Impl()
	user, err := userImpl.Get(ctx, loginId)
	if err != nil {
		return errs.Sys(err)
	}
	if user == nil {
		return errs.ERR_NOT_ACCOUNT
	}

	upd := map[string]any{}
	if nickName != "" {
		upd["nick_name"] = nickName
	}
	if email != "" {
		upd["email"] = email
	}
	if phoneArea != "" {
		upd["phone_area"] = phoneArea
	}
	if phone != "" {
		upd["phone"] = phone
	}
	if sex >= 1 && sex <= 3 {
		upd["sex"] = sex
	}
	if address != "" {
		upd["address"] = address
	}
	if len(upd) == 0 {
		return nil
	}

	_, err = m.engine.Context(ctx).Table(&model.SysUser{}).ID(user.Id).Update(upd)
	if err != nil {
		return errs.Sys(err)
	}
	// 同步更新 session，使 SSE 等能读到最新昵称
	if nickName != "" {
		_ = auth.SetSess(ctx, loginId, map[string]any{"nickName": nickName})
	}
	return nil
}
