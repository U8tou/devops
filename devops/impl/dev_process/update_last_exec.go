package devprocess

import (
	"context"
	"devops/model"
)

func (m *DevProcessImpl) UpdateLastExec(ctx context.Context, id int64, lastExecTime int64, lastExecDurationMs int64, lastExecResult string, lastExecLog string, updateBy int64) (int64, error) {
	// MustCols：避免零值/omit 规则导致 last_exec_result、last_exec_log、last_exec_duration_ms 未写入
	return m.engine.Context(ctx).ID(id).MustCols("last_exec_time", "last_exec_duration_ms", "last_exec_result", "last_exec_log", "update_by").Update(&model.DevProcess{
		LastExecTime:       lastExecTime,
		LastExecDurationMs: lastExecDurationMs,
		LastExecResult:     lastExecResult,
		LastExecLog:        lastExecLog,
		UpdateBy:           updateBy,
	})
}
