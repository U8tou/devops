package syspost

import (
	"context"
	"system/model"
)

// List 分页
func (m *SysPostImpl) List(ctx context.Context, args *model.SysPostDto) (int64, []model.SysPostVo, error) {
	query := m.engine.Context(ctx).Table(&model.SysPost{}).Where("1 = 1")
	if args.Id != 0 {
		query.And("id = ?", args.Id)
	}
	if args.Name != "" {
		query.And("name LIKE ?", "%"+args.Name+"%")
	}
	if args.Status != 0 {
		query.And("status = ?", args.Status)
	}

	var total int64
	var err error
	rows := make([]model.SysPostVo, 0)
	if args.Size > 0 {
		total, err = query.Limit(args.Size, (args.Current-1)*args.Size).Asc("sort", "id").FindAndCount(&rows)
	} else {
		err = query.Asc("sort", "id").Find(&rows)
		if err == nil {
			total = int64(len(rows))
		}
	}
	return total, rows, err
}
