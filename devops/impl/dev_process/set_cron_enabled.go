package devprocess

import (
	"context"
	"devops/model"
	"strings"
	"time"

	"pkg/cronvalidate"
	"pkg/errs"
)

// SetCronEnabled 更新定时启用；开启时若未配置表达式则写入默认每分钟。
func (m *DevProcessImpl) SetCronEnabled(ctx context.Context, id int64, enabled int8, updateBy int64) (int64, error) {
	if enabled != 0 && enabled != 1 {
		return 0, errs.New("启用状态无效")
	}
	var row model.DevProcess
	has, err := m.engine.Context(ctx).ID(id).Get(&row)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errs.ERR_DB_NO_EXIST
	}
	now := time.Now().Unix()
	if enabled == 1 {
		expr := strings.TrimSpace(row.CronExpr)
		if expr == "" {
			expr = DefaultCronExprEveryMinute
		}
		if err := cronvalidate.ValidateExpr(expr); err != nil {
			return 0, errs.New("Cron 表达式无效")
		}
		affected, err := m.engine.Context(ctx).ID(id).Cols("cron_enabled", "cron_expr", "update_by", "update_time").Update(&model.DevProcess{
			CronEnabled: 1,
			CronExpr:    expr,
			UpdateBy:    updateBy,
			UpdateTime:  now,
		})
		return affected, err
	}
	affected, err := m.engine.Context(ctx).ID(id).Cols("cron_enabled", "update_by", "update_time").Update(&model.DevProcess{
		CronEnabled: 0,
		UpdateBy:    updateBy,
		UpdateTime:  now,
	})
	return affected, err
}
