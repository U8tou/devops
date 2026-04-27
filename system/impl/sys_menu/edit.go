package sysmenu

import (
	"context"
	"system/model"
)

// Edit 编辑
func (m *SysMenuImpl) Edit(ctx context.Context, args *model.SysMenuDto) (int64, error) {
	menu := args.SysMenu
	return m.engine.Context(ctx).Where("id = ?", menu.Id).
		Cols("pid", "types", "permis", "remark", "update_by", "update_time").
		Update(&menu)
}
