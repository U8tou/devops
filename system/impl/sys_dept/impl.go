package sysdept

import (
	"context"
	"pkg/db"
	"sync"
	"system/model"

	"xorm.io/xorm"
)

/**
Notes: SysDept Impl
Time:  2025-04-29 10:48:59
*/

var (
	sysDeptImpl     *SysDeptImpl
	sysDeptImplOnce sync.Once
)

// SysDeptImpl 类
type SysDeptImpl struct {
	engine *xorm.Engine
}

// Impl 实例化 | 依赖注入
func Impl() ISysDeptImpl {
	sysDeptImplOnce.Do(func() {
		sysDeptImpl = &SysDeptImpl{
			engine: db.GetDb(),
		}
	})
	return sysDeptImpl
}

// ISysDeptImpl 接口
type ISysDeptImpl interface {
	// List 分页
	List(ctx context.Context, args *model.SysDeptDto) (int64, []model.SysDeptVo, error)
	// Get 查询
	Get(ctx context.Context, id string) (*model.SysDeptVo, error)
	// Del 删除
	Del(ctx context.Context, ids []int64) (int64, error)
	// Add 新增
	Add(ctx context.Context, args *model.SysDeptDto) (int64, error)
	// Edit 编辑
	Edit(ctx context.Context, args *model.SysDeptDto) (int64, error)
}
