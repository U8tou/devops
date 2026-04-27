package sysauth

import (
	"context"
	"log"
	"pkg/auth"
	"pkg/crypto"
	"pkg/errs"
	"pkg/tools/datacv"
	"system/model"

	sysmenu "system/impl/sys_menu"
	sysuser "system/impl/sys_user"
)

// Login 登录：从 DB 查用户，bcrypt 校验密码，设置权限与 Token
func (_ *SysAuthImpl) Login(ctx context.Context, args *model.LoginReq) (*model.LoginResp, error) {
	isLogin := false
	_, err := auth.TryLogin(ctx, args.UserName)
	if err != nil {
		return nil, err
	}
	defer func() {
		if isLogin {
			_ = auth.TryClean(ctx, args.UserName)
		}
	}()

	// 按用户名查用户
	userImpl := sysuser.Impl()
	user, err := userImpl.GetByUserName(ctx, args.UserName)
	if err != nil {
		return nil, errs.Sys(err)
	}
	if user == nil {
		return nil, errs.ERR_NOT_ACCOUNT
	}

	// 校验密码（支持 bcrypt 与明文兼容）
	if !crypto.CheckPasswordHash(args.Password, user.Password) {
		return nil, errs.ERR_PWD
	}

	// 状态检查
	if user.Status != 1 {
		return nil, errs.ERR_SYS_405
	}

	// 根用户（userType=00）：全量菜单；普通用户须绑定角色
	menuImpl := sysmenu.Impl()
	var menus []model.SysMenuVo
	fullAccess := user.UserType == model.UserTypeRoot
	if fullAccess {
		menus, _ = menuImpl.All(ctx)
	} else {
		roleIds, _ := userImpl.GetUserRoleIds(ctx, user.Id)
		if len(roleIds) == 0 {
			return nil, errs.ERR_NO_ROLE
		}
		menus, _ = menuImpl.GetByRoleIds(ctx, roleIds)
	}

	ms := make([]string, 0)
	ps := make([]string, 0)
	for _, m := range menus {
		if m.Permis != "" {
			if m.Types == 1 {
				// 菜单权限
				ms = append(ms, m.Permis)
			} else {
				// 操作权限
				ps = append(ps, m.Permis)
			}
		}
	}
	// ms = append(ms, "R_SUPER")
	// ps = append(ps, "B_CODE1", "B_CODE2", "B_CODE3", "sys:user:edit", "sys:user:delete")

	loginId := datacv.IntToStr(user.Id)
	if err = auth.SetMenu(ctx, loginId, ms); err != nil {
		return nil, errs.Sys(err)
	}
	if err = auth.SetPermis(ctx, loginId, ps); err != nil {
		return nil, errs.Sys(err)
	}

	// 设置会话
	mp := make(map[string]any, 6)
	mp["loginId"] = user.Id
	mp["userName"] = user.UserName
	mp["nickName"] = user.NickName
	mp["avatar"] = user.Avatar
	mp["email"] = user.Email
	mp["userType"] = user.UserType
	err = auth.SetSess(ctx, loginId, mp)
	if err != nil {
		return nil, err
	}

	device := "web"
	tokenInfo := auth.GetToken(ctx, loginId, device)
	log.Println("✅ Login ok! userName:", args.UserName)

	isLogin = true
	return &model.LoginResp{
		Token:        tokenInfo.AccessToken,
		RefreshToken: tokenInfo.RefreshToken,
	}, nil
}
