package sysmenu

import (
	"context"
	"system/model"
)

// Get 查询
func (m *SysMenuImpl) Get(ctx context.Context, id int64) (*model.SysMenuVo, error) {
	var row model.SysMenuVo
	has, err := m.engine.Context(ctx).Table(&model.SysMenu{}).Where("id = ?", id).Get(&row)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &row, nil
}
