package sysrole

import (
	"context"
	"system/model"
)

// Get 查询
func (m *SysRoleImpl) Get(ctx context.Context, id int64) (*model.SysRoleVo, error) {
	var row model.SysRoleVo
	has, err := m.engine.Context(ctx).Table(&model.SysRole{}).Where("id = ?", id).Get(&row)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	// 获取角色的菜单ID列表
	menuIds, err := m.GetRoleMenuIds(ctx, id)
	if err != nil {
		return nil, err
	}
	row.Menus = menuIds
	// 获取角色的部门ID列表
	deptIds, err := m.GetRoleDeptIds(ctx, id)
	if err != nil {
		return nil, err
	}
	row.Depts = deptIds
	return &row, nil
}

// GetRoleMenuIds 获取角色的菜单ID列表
func (m *SysRoleImpl) GetRoleMenuIds(ctx context.Context, roleId int64) ([]int64, error) {
	var menuIds []int64
	err := m.engine.Context(ctx).Table(&model.SysRoleMenu{}).
		Where("role_id = ?", roleId).
		Cols("menu_id").
		Find(&menuIds)
	return menuIds, err
}

// GetRoleDeptIds 获取角色的部门ID列表
func (m *SysRoleImpl) GetRoleDeptIds(ctx context.Context, roleId int64) ([]int64, error) {
	var deptIds []int64
	err := m.engine.Context(ctx).Table(&model.SysRoleDept{}).
		Where("role_id = ?", roleId).
		Cols("dept_id").
		Find(&deptIds)
	return deptIds, err
}
