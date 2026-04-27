package devprocess

import (
	"context"
	"devops/model"
	"strings"
	"time"

	"pkg/errs"

	"xorm.io/xorm"
)

func (m *DevProcessImpl) Edit(ctx context.Context, args *model.DevProcessDto) (int64, error) {
	v := args.DevProcess
	if v.CronEnabled != 0 && v.CronEnabled != 1 {
		v.CronEnabled = 0
	}
	v.CronExpr = strings.TrimSpace(v.CronExpr)
	if v.CronEnabled == 1 && v.CronExpr == "" {
		v.CronExpr = DefaultCronExprEveryMinute
	}
	if err := validateCronSettings(v.CronEnabled, v.CronExpr); err != nil {
		return 0, err
	}
	var old model.DevProcess
	has, err := m.engine.Context(ctx).ID(v.Id).Get(&old)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errs.ERR_DB_NO_EXIST
	}
	v.UpdateTime = time.Now().Unix()
	if len(args.TagIds) > 0 {
		if err := m.validateTagIdsExist(ctx, args.TagIds); err != nil {
			return 0, err
		}
	}
	var affect int64
	_, err = m.engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		s = s.Context(ctx)
		a, err := s.ID(v.Id).Cols("code", "remark", "cron_expr", "cron_enabled", "update_by", "update_time").Update(&v)
		if err != nil {
			return nil, err
		}
		affect = a
		err = replaceTagLinksSession(s, v.Id, dedupeTagIds(args.TagIds))
		return nil, err
	})
	return affect, err
}
