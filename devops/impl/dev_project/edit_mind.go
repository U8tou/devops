package devproject

import (
	"context"
	"devops/model"
	"time"

	"pkg/errs"
)

func (m *DevProjectImpl) EditMind(ctx context.Context, id int64, mindJson string, updateBy int64) (int64, error) {
	var old model.DevProject
	has, err := m.engine.Context(ctx).ID(id).Get(&old)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errs.ERR_DB_NO_EXIST
	}
	up := model.DevProject{
		MindJson:   mindJson,
		UpdateBy:   updateBy,
		UpdateTime: time.Now().Unix(),
	}
	return m.engine.Context(ctx).ID(id).Cols("mind_json", "update_by", "update_time").Update(&up)
}
