package sysrole

import (
	"context"
	"system/model"
)

// AssignDept 分配部门权限
func (m *SysRoleImpl) AssignDept(ctx context.Context, roleId int64, deptIds []int64, deptLinkage int8) (int64, error) {
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 更新角色是否联动菜单
	role := model.SysRole{
		DeptLinkage: deptLinkage,
	}
	_, err := session.Where("id = ?", roleId).
		Cols("dept_linkage", "update_time").
		Update(&role)
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 删除旧的角色部门关联
	affect, err := session.Where("role_id = ?", roleId).Delete(&model.SysRoleDept{})
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 新增新的角色部门关联
	if len(deptIds) > 0 {
		roleDepts := make([]model.SysRoleDept, len(deptIds))
		for i, deptId := range deptIds {
			roleDepts[i] = model.SysRoleDept{
				RoleId: roleId,
				DeptId: deptId,
			}
		}
		_, err = session.Insert(&roleDepts)
		if err != nil {
			_ = session.Rollback()
			return 0, err
		}
	}

	return affect, session.Commit()
}
