package syspost

import (
	"context"
	"system/model"
)

// Get 查
func (m *SysPostImpl) Get(ctx context.Context, id string) (*model.SysPostVo, error) {
	var v model.SysPostVo
	has, err := m.engine.Context(ctx).Table(&model.SysPost{}).ID(id).Get(&v)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &v, nil
}
