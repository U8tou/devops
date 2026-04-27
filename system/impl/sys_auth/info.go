package sysauth

import (
	"context"
	"fmt"
	"pkg/auth"
	"pkg/conf"
	"pkg/tools/datacv"
	"system/model"

	sysuser "system/impl/sys_user"
)

// Info 登录信息
func (_ *SysAuthImpl) Info(ctx context.Context, loginId string) (*model.InfoResp, error) {
	// 从 DB 获取完整用户信息
	userImpl := sysuser.Impl()
	user, err := userImpl.Get(ctx, loginId)
	if err != nil {
		return nil, fmt.Errorf("sys_user.Get(ctx, %s): %w", loginId, err)
	}
	if user == nil {
		return nil, fmt.Errorf("sys_user.Get(ctx, %s): user not found", loginId)
	}

	menus, err := auth.GetMenu(ctx, loginId)
	if err != nil {
		return nil, fmt.Errorf("auth.GetMenu(ctx, loginId): %w", err)
	}
	permis, err := auth.GetPermis(ctx, loginId)
	if err != nil {
		return nil, fmt.Errorf("auth.GetPermis(ctx, loginId): %w", err)
	}
	avatar := conf.FileUrl(user.Avatar)

	return &model.InfoResp{
		UserId:     user.Id,
		UserName:   user.UserName,
		NickName:   user.NickName,
		UserType:   user.UserType,
		Email:      user.Email,
		PhoneArea:  user.PhoneArea,
		Phone:      user.Phone,
		Sex:        user.Sex,
		Avatar:     avatar,
		Status:     user.Status,
		Address:    user.Address,
		Remark:     user.Remark,
		CreateTime: datacv.IntToStr(user.CreateTime),
		CreateBy:   datacv.IntToStr(user.CreateBy),
		UpdateTime: datacv.IntToStr(user.UpdateTime),
		UpdateBy:   datacv.IntToStr(user.UpdateBy),
		Depts:      user.Depts,
		Roles:      user.Roles,
		Posts:      user.Posts,
		Buttons:    permis,
		Menus:      menus,
	}, nil
}
