package sysdept

import (
	"context"
	"system/model"
)

// Get 查
func (m *SysDeptImpl) Get(ctx context.Context, id string) (*model.SysDeptVo, error) {
	var v model.SysDeptVo
	has, err := m.engine.Context(ctx).Table(&model.SysDept{}).ID(id).Get(&v)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &v, nil
}
