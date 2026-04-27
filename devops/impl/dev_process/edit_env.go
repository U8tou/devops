package devprocess

import (
	"context"
	"devops/model"
	"time"
)

func (m *DevProcessImpl) EditEnv(ctx context.Context, id int64, envJson string, updateBy int64) (int64, error) {
	v := model.DevProcess{
		EnvJson:  envJson,
		UpdateBy: updateBy,
	}
	v.UpdateTime = time.Now().Unix()
	return m.engine.Context(ctx).ID(id).Cols("env_json", "update_by", "update_time").Update(&v)
}
