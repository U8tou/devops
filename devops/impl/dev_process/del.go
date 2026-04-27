package devprocess

import (
	"context"
	"devops/model"
)

// Del 软删（xorm 对带 `deleted` 标记的 DeleteTime 字段执行更新，写入删除时间戳，不物理删行）；并清理 dev_process_tag_link。
func (m *DevProcessImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	if err := m.DeleteLinksByProcessIds(ctx, ids); err != nil {
		return 0, err
	}
	return m.engine.Context(ctx).In("id", ids).Delete(&model.DevProcess{})
}
