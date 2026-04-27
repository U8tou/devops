package devproject

import (
	"context"
	"devops/model"
)

// Del 软删；并清理 dev_project_tag_link。
func (m *DevProjectImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	if err := m.DeleteLinksByProjectIds(ctx, ids); err != nil {
		return 0, err
	}
	return m.engine.Context(ctx).In("id", ids).Delete(&model.DevProject{})
}
