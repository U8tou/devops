package sysuser

import (
	"context"
	"pkg/crypto"
	"system/model"
)

// Edit 改（密码为空表示不修改密码）
func (m *SysUserImpl) Edit(ctx context.Context, args *model.SysUserDto) (int64, error) {
	v := args.SysUser
	if args.Password != "" {
		hashedPwd, err := crypto.HashPassword(args.Password)
		if err != nil {
			return 0, err
		}
		v.Password = hashedPwd
		return m.engine.Context(ctx).Table(&model.SysUser{}).ID(v.Id).AllCols().Update(&v)
	}
	return m.engine.Context(ctx).Table(&model.SysUser{}).ID(v.Id).Omit("password").Update(&v)
}
