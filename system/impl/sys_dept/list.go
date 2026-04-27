package sysdept

import (
	"context"
	"system/model"
)

// List 分页
func (m *SysDeptImpl) List(ctx context.Context, args *model.SysDeptDto) (int64, []model.SysDeptVo, error) {
	// 构建筛选条件（排除软删除）
	query := m.engine.Context(ctx).Table(&model.SysDept{}).Where("1 = 1")
	if args.Email != "" {
		query.And("email = ?", args.Email)
	}
	if args.Leader != "" {
		query.And("leader = ?", args.Leader)
	}
	if args.Name != "" {
		query.And("name = ?", args.Name)
	}
	if args.Phone != "" {
		query.And("phone = ?", args.Phone)
	}
	if args.Profile != "" {
		query.And("profile = ?", args.Profile)
	}
	if args.Status != 0 {
		query.And("status = ?", args.Status)
	}
	// 查询
	var total int64
	var err error
	rows := make([]model.SysDeptVo, 0)
	if args.Size > 0 {
		total, err = query.Limit(args.Size, (args.Current-1)*args.Size).Asc("sort", "id").FindAndCount(&rows)
	} else {
		err = query.Asc("id").Find(&rows)
	}
	return total, rows, err
}
