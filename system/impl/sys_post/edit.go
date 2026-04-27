package syspost

import (
	"context"
	"system/model"
)

// Edit 改
func (m *SysPostImpl) Edit(ctx context.Context, args *model.SysPostDto) (int64, error) {
	v := args.SysPost
	return m.engine.Context(ctx).ID(v.Id).Update(&v)
}
