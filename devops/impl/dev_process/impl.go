package devprocess

import (
	"context"
	"devops/model"
	"pkg/db"
	"sync"

	"xorm.io/xorm"
)

var (
	devProcessImpl     *DevProcessImpl
	devProcessImplOnce sync.Once
)

type DevProcessImpl struct {
	engine *xorm.Engine
}

func Impl() IDevProcessImpl {
	devProcessImplOnce.Do(func() {
		devProcessImpl = &DevProcessImpl{
			engine: db.GetDb(),
		}
	})
	return devProcessImpl
}

type IDevProcessImpl interface {
	List(ctx context.Context, args *model.DevProcessDto) (int64, []model.DevProcessVo, error)
	Get(ctx context.Context, id string) (*model.DevProcessVo, error)
	Add(ctx context.Context, args *model.DevProcessDto) (int64, error)
	Edit(ctx context.Context, args *model.DevProcessDto) (int64, error)
	EditFlow(ctx context.Context, args *model.DevProcessDto) (int64, error)
	EditEnv(ctx context.Context, id int64, envJson string, updateBy int64) (int64, error)
	SetCronEnabled(ctx context.Context, id int64, enabled int8, updateBy int64) (int64, error)
	UpdateLastExec(ctx context.Context, id int64, lastExecTime int64, lastExecDurationMs int64, lastExecResult string, lastExecLog string, updateBy int64) (int64, error)
	Del(ctx context.Context, ids []int64) (int64, error)
	// ListCronScheduled 返回已启用定时的流程（用于启动全局 Cron 调度）
	ListCronScheduled(ctx context.Context) ([]model.DevCronScheduleRow, error)
	// 标签字典
	TagList(ctx context.Context) ([]model.DevProcessTag, error)
	TagAdd(ctx context.Context, name string) (int64, error)
	TagEdit(ctx context.Context, id int64, name string) (int64, error)
	TagDel(ctx context.Context, id int64) (int64, error)
}
