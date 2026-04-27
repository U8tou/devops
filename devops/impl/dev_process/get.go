package devprocess

import (
	"context"
	"devops/model"
)

func (m *DevProcessImpl) Get(ctx context.Context, id string) (*model.DevProcessVo, error) {
	var v model.DevProcessVo
	has, err := m.engine.Context(ctx).Table(&model.DevProcess{}).ID(id).Get(&v)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	tags, err := m.LoadProcessTags(ctx, v.Id)
	if err != nil {
		return nil, err
	}
	v.Tags = tags
	return &v, nil
}
