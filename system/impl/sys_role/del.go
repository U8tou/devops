package sysrole

import (
	"context"
	"system/model"
)

// Del 软删（使用 xorm deleted，关联表同步清理）
func (m *SysRoleImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 软删角色（xorm 自动设置 delete_time）
	affect, err := session.In("id", ids).Delete(&model.SysRole{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	// 删除角色菜单/部门/用户关联（硬删，避免悬空引用）
	_, err = session.In("role_id", ids).Delete(&model.SysRoleMenu{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	// 删除角色部门关联
	_, err = session.In("role_id", ids).Delete(&model.SysRoleDept{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	// 删除用户角色关联
	_, err = session.In("role_id", ids).Delete(&model.SysUserRole{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	return affect, session.Commit()
}
