package devproject

import (
	"context"
	"devops/model"
)

func (m *DevProjectImpl) Get(ctx context.Context, id string) (*model.DevProjectVo, error) {
	var v model.DevProjectVo
	has, err := m.engine.Context(ctx).Table(&model.DevProject{}).ID(id).Get(&v)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	tags, err := m.LoadProjectTags(ctx, v.Id)
	if err != nil {
		return nil, err
	}
	v.Tags = tags
	return &v, nil
}
