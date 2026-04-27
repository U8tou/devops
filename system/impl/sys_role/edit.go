package sysrole

import (
	"context"
	"system/model"
)

// Edit 编辑
func (m *SysRoleImpl) Edit(ctx context.Context, args *model.SysRoleDto) (int64, error) {
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	role := args.SysRole
	affect, err := session.Where("id = ?", role.Id).
		Cols("name", "role", "status", "menu_linkage", "dept_linkage", "sort", "remark", "update_by", "update_time").
		Update(&role)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	return affect, session.Commit()
}
