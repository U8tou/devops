package sysmenu

import (
	"context"
	"system/model"
)

// Del 删除
func (m *SysMenuImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 删除菜单
	affect, err := session.In("id", ids).Delete(&model.SysMenu{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	// 删除角色菜单关联
	_, err = session.In("menu_id", ids).Delete(&model.SysRoleMenu{})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	return affect, session.Commit()
}
