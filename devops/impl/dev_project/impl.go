package devproject

import (
	"context"
	"devops/model"
	"pkg/db"
	"sync"

	"xorm.io/xorm"
)

var (
	devProjectImpl     *DevProjectImpl
	devProjectImplOnce sync.Once
)

type DevProjectImpl struct {
	engine *xorm.Engine
}

func Impl() IDevProjectImpl {
	devProjectImplOnce.Do(func() {
		devProjectImpl = &DevProjectImpl{
			engine: db.GetDb(),
		}
	})
	return devProjectImpl
}

type IDevProjectImpl interface {
	List(ctx context.Context, args *model.DevProjectDto) (int64, []model.DevProjectVo, error)
	Get(ctx context.Context, id string) (*model.DevProjectVo, error)
	Add(ctx context.Context, args *model.DevProjectDto) (int64, error)
	Edit(ctx context.Context, args *model.DevProjectDto) (int64, error)
	EditMind(ctx context.Context, id int64, mindJson string, updateBy int64) (int64, error)
	Del(ctx context.Context, ids []int64) (int64, error)
	TagList(ctx context.Context) ([]model.DevProjectTag, error)
	TagAdd(ctx context.Context, name string) (int64, error)
	TagEdit(ctx context.Context, id int64, name string) (int64, error)
	TagDel(ctx context.Context, id int64) (int64, error)
}
