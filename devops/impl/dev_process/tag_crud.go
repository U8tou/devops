package devprocess

import (
	"context"
	"devops/model"
	"strings"

	"github.com/yitter/idgenerator-go/idgen"
)

// TagList 返回全部标签（字典顺序）
func (m *DevProcessImpl) TagList(ctx context.Context) ([]model.DevProcessTag, error) {
	var rows []model.DevProcessTag
	err := m.engine.Context(ctx).Table(&model.DevProcessTag{}).Asc("name").Find(&rows)
	return rows, err
}

// TagAdd 新增标签
func (m *DevProcessImpl) TagAdd(ctx context.Context, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errTagNameEmpty
	}
	row := model.DevProcessTag{Id: idgen.NextId(), Name: name}
	_, err := m.engine.Context(ctx).InsertOne(&row)
	if err != nil {
		return 0, err
	}
	return row.Id, nil
}

// TagEdit 更新标签名
func (m *DevProcessImpl) TagEdit(ctx context.Context, id int64, name string) (int64, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return 0, errTagNameEmpty
	}
	return m.engine.Context(ctx).ID(id).Cols("name").Update(&model.DevProcessTag{Name: name})
}

// TagDel 仅删除字典行，不删 dev_process_tag_link（保留孤儿关联）
func (m *DevProcessImpl) TagDel(ctx context.Context, id int64) (int64, error) {
	return m.engine.Context(ctx).ID(id).Delete(&model.DevProcessTag{})
}
