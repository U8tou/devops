package devprocess

import (
	"context"
	"devops/model"
	"time"
)

func (m *DevProcessImpl) EditFlow(ctx context.Context, args *model.DevProcessDto) (int64, error) {
	v := args.DevProcess
	v.UpdateTime = time.Now().Unix()
	return m.engine.Context(ctx).ID(v.Id).Cols("flow", "update_by", "update_time").Update(&v)
}
