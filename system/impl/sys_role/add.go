package sysrole

import (
	"context"
	"system/model"

	"github.com/yitter/idgenerator-go/idgen"
)

// Add 新增
func (m *SysRoleImpl) Add(ctx context.Context, args *model.SysRoleDto) (int64, error) {
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	role := args.SysRole
	role.Id = idgen.NextId()
	role.DeleteTime = 0
	if _, err := session.Insert(&role); err != nil {
		session.Rollback()
		return 0, err
	}
	return role.Id, session.Commit()
}
