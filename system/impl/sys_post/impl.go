package syspost

import (
	"context"
	"pkg/db"
	"sync"
	"system/model"

	"xorm.io/xorm"
)

var (
	sysPostImpl     *SysPostImpl
	sysPostImplOnce sync.Once
)

// SysPostImpl 类
type SysPostImpl struct {
	engine *xorm.Engine
}

// Impl 实例化
func Impl() ISysPostImpl {
	sysPostImplOnce.Do(func() {
		sysPostImpl = &SysPostImpl{
			engine: db.GetDb(),
		}
	})
	return sysPostImpl
}

// ISysPostImpl 接口
type ISysPostImpl interface {
	List(ctx context.Context, args *model.SysPostDto) (int64, []model.SysPostVo, error)
	Get(ctx context.Context, id string) (*model.SysPostVo, error)
	Del(ctx context.Context, ids []int64) (int64, error)
	Add(ctx context.Context, args *model.SysPostDto) (int64, error)
	Edit(ctx context.Context, args *model.SysPostDto) (int64, error)
}
