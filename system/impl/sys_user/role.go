package sysuser

import (
	"context"
	"system/model"
)

// AssignRole 分配角色
func (m *SysUserImpl) AssignRole(ctx context.Context, userId int64, roleIds []int64) (int64, error) {
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 删除旧的用户角色关联
	affect, err := session.Where("user_id = ?", userId).Delete(&model.SysUserRole{})
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 新增新的用户角色关联
	if len(roleIds) > 0 {
		userRoles := make([]model.SysUserRole, len(roleIds))
		for i, roleId := range roleIds {
			userRoles[i] = model.SysUserRole{
				UserId: userId,
				RoleId: roleId,
			}
		}
		_, err = session.Insert(&userRoles)
		if err != nil {
			session.Rollback()
			return 0, err
		}
	}

	return affect, session.Commit()
}

// GetUserRoleIds 获取用户角色ID列表
func (m *SysUserImpl) GetUserRoleIds(ctx context.Context, userId int64) ([]int64, error) {
	var roleIds []int64
	err := m.engine.Context(ctx).Table(&model.SysUserRole{}).
		Where("user_id = ?", userId).
		Cols("role_id").
		Find(&roleIds)
	return roleIds, err
}
