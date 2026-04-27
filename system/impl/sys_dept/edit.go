package sysdept

import (
	"context"
	"system/model"
)

// Edit 改
func (m *SysDeptImpl) Edit(ctx context.Context, args *model.SysDeptDto) (int64, error) {
	v := args.SysDept
	return m.engine.Context(ctx).ID(v.Id).Update(&v)
}
