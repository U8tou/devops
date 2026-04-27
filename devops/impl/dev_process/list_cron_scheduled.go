package devprocess

import (
	"context"
	"devops/model"
)

// ListCronScheduled 查询「定时执行」且「定时已启用」的流程，用于进程内 Cron 注册。
func (m *DevProcessImpl) ListCronScheduled(ctx context.Context) ([]model.DevCronScheduleRow, error) {
	var rows []model.DevCronScheduleRow
	err := m.engine.Context(ctx).Table(&model.DevProcess{}).
		Where("cron_enabled = ?", 1).
		Cols("id", "cron_expr").
		Find(&rows)
	return rows, err
}
