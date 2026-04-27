package sysrole

import (
	"context"
	"system/model"
)

// List 分页
func (m *SysRoleImpl) List(ctx context.Context, args *model.SysRoleDto) (int64, []model.SysRoleVo, error) {
	// 构建筛选条件
	query := m.engine.Context(ctx).Table(&model.SysRole{}).Where("1 = 1")
	if args.Id != 0 {
		query.And("id = ?", args.Id)
	}
	if args.Name != "" {
		query.And("name like ?", "%"+args.Name+"%")
	}
	if args.Role != "" {
		query.And("role like ?", "%"+args.Role+"%")
	}
	if args.Status != 0 {
		query.And("status = ?", args.Status)
	}

	// 查询
	var total int64
	var err error
	rows := make([]model.SysRoleVo, 0)
	if args.Size > 0 {
		total, err = query.Limit(args.Size, (args.Current-1)*args.Size).Asc("sort").FindAndCount(&rows)
	} else {
		err = query.Asc("sort").Find(&rows)
	}
	return total, rows, err
}

// All 全量
func (m *SysRoleImpl) All(ctx context.Context) ([]model.SysRoleVo, error) {
	rows := make([]model.SysRoleVo, 0)
	err := m.engine.Context(ctx).Table(&model.SysRole{}).Where("status = 1").Asc("sort").Find(&rows)
	return rows, err
}
