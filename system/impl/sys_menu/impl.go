package sysmenu

import (
	"context"
	"pkg/db"
	"sync"
	"system/model"

	"xorm.io/xorm"
)

/**
Notes: 系统菜单/权限 IMPL
Time:  2025-01-29
*/

var (
	sysMenuImpl     *SysMenuImpl
	sysMenuImplOnce sync.Once
)

// SysMenuImpl 类
type SysMenuImpl struct {
	engine *xorm.Engine
}

// Impl 实例化 | 依赖注入
func Impl() ISysMenuImpl {
	sysMenuImplOnce.Do(func() {
		sysMenuImpl = &SysMenuImpl{
			engine: db.GetDb(),
		}
	})
	return sysMenuImpl
}

// ISysMenuImpl 接口
type ISysMenuImpl interface {
	// List 分页
	List(ctx context.Context, args *model.SysMenuDto) (int64, []model.SysMenuVo, error)
	// Get 查询
	Get(ctx context.Context, id int64) (*model.SysMenuVo, error)
	// Del 删除
	Del(ctx context.Context, ids []int64) (int64, error)
	// Add 新增
	Add(ctx context.Context, args *model.SysMenuDto) (int64, error)
	// Edit 编辑
	Edit(ctx context.Context, args *model.SysMenuDto) (int64, error)
	// All 全量
	All(ctx context.Context) ([]model.SysMenuVo, error)
	// GetByRoleIds 根据角色ID列表获取菜单权限
	GetByRoleIds(ctx context.Context, roleIds []int64) ([]model.SysMenuVo, error)
}
