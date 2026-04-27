package sysrole

import (
	"context"
	"system/model"
)

// AssignMenu 分配菜单权限
func (m *SysRoleImpl) AssignMenu(ctx context.Context, roleId int64, menuIds []int64, menuLinkage int8) (int64, error) {
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 更新角色是否联动菜单
	role := model.SysRole{
		MenuLinkage: menuLinkage,
	}
	_, err := session.Where("id = ?", roleId).
		Cols("menu_linkage", "update_time").
		Update(&role)
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 删除旧的角色菜单关联
	affect, err := session.Where("role_id = ?", roleId).Delete(&model.SysRoleMenu{})
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 新增新的角色菜单关联
	if len(menuIds) > 0 {
		roleMenus := make([]model.SysRoleMenu, len(menuIds))
		for i, menuId := range menuIds {
			roleMenus[i] = model.SysRoleMenu{
				RoleId: roleId,
				MenuId: menuId,
			}
		}
		_, err = session.Insert(&roleMenus)
		if err != nil {
			session.Rollback()
			return 0, err
		}
	}

	return affect, session.Commit()
}
