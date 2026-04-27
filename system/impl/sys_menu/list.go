package sysmenu

import (
	"context"
	"system/model"
)

// List 分页
func (m *SysMenuImpl) List(ctx context.Context, args *model.SysMenuDto) (int64, []model.SysMenuVo, error) {
	// 构建筛选条件
	query := m.engine.Context(ctx).Table(&model.SysMenu{}).Where("1 = 1")
	if args.Id != 0 {
		query.And("id = ?", args.Id)
	}
	if args.Pid != 0 {
		query.And("pid = ?", args.Pid)
	}
	if args.Types != 0 {
		query.And("types = ?", args.Types)
	}
	if args.Permis != "" {
		query.And("permis like ?", "%"+args.Permis+"%")
	}

	// 查询
	var total int64
	var err error
	rows := make([]model.SysMenuVo, 0)
	if args.Size > 0 {
		total, err = query.Limit(args.Size, (args.Current-1)*args.Size).Asc("id").FindAndCount(&rows)
	} else {
		err = query.Asc("id").Find(&rows)
	}
	return total, rows, err
}

// All 全量
func (m *SysMenuImpl) All(ctx context.Context) ([]model.SysMenuVo, error) {
	rows := make([]model.SysMenuVo, 0)
	err := m.engine.Context(ctx).Table(&model.SysMenu{}).Asc("id").Find(&rows)
	return rows, err
}

// GetByRoleIds 根据角色ID列表获取菜单权限
func (m *SysMenuImpl) GetByRoleIds(ctx context.Context, roleIds []int64) ([]model.SysMenuVo, error) {
	if len(roleIds) == 0 {
		return []model.SysMenuVo{}, nil
	}
	rows := make([]model.SysMenuVo, 0)
	err := m.engine.Context(ctx).Table(&model.SysMenu{}).
		Join("INNER", &model.SysRoleMenu{}, "sys_menu.id = sys_role_menu.menu_id").
		In("sys_role_menu.role_id", roleIds).
		Distinct("sys_menu.*").
		Find(&rows)
	return rows, err
}
