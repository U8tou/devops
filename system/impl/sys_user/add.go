package sysuser

import (
	"context"
	"pkg/crypto"
	"system/model"

	"github.com/yitter/idgenerator-go/idgen"
)

// Add 增
func (m *SysUserImpl) Add(ctx context.Context, args *model.SysUserDto) (int64, error) {
	hashedPwd, err := crypto.HashPassword(args.Password)
	if err != nil {
		return 0, err
	}
	// DTO 内嵌了 SysUser，先整体拷贝再只覆盖入库时需要定制的字段
	v := args.SysUser
	v.Id = idgen.NextId()
	v.DeleteTime = 0
	v.Password = hashedPwd
	return m.engine.Context(ctx).InsertOne(&v)
}
