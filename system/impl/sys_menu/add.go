package sysmenu

import (
	"context"
	"system/model"

	"github.com/yitter/idgenerator-go/idgen"
)

// Add 新增
func (m *SysMenuImpl) Add(ctx context.Context, args *model.SysMenuDto) (int64, error) {
	menu := args.SysMenu
	menu.Id = idgen.NextId()
	_, err := m.engine.Context(ctx).Insert(&menu)
	if err != nil {
		return 0, err
	}
	return menu.Id, nil
}
