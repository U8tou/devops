package sysdept

import (
	"context"
	"system/model"
)

// Del 软删（xorm deleted 自动设置 delete_time 为秒级时间戳）
func (m *SysDeptImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	return m.engine.Context(ctx).In("id", ids).Delete(&model.SysDept{})
}
