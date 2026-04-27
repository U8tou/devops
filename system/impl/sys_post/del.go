package syspost

import (
	"context"
	"system/model"
)

// Del 软删
func (m *SysPostImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return m.engine.Context(ctx).In("id", ids).Delete(&model.SysPost{})
}
