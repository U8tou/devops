package devprocess

import (
	"context"
	"devops/model"
	"strings"

	"github.com/yitter/idgenerator-go/idgen"
	"xorm.io/xorm"
)

func (m *DevProcessImpl) Add(ctx context.Context, args *model.DevProcessDto) (int64, error) {
	if args.CronEnabled != 0 && args.CronEnabled != 1 {
		args.CronEnabled = 0
	}
	args.CronExpr = strings.TrimSpace(args.CronExpr)
	if args.CronEnabled == 1 && args.CronExpr == "" {
		args.CronExpr = DefaultCronExprEveryMinute
	}
	if err := validateCronSettings(args.CronEnabled, args.CronExpr); err != nil {
		return 0, err
	}
	args.Id = idgen.NextId()
	args.DeleteTime = 0
	args.LastExecTime = 0
	args.LastExecResult = ""
	args.LastExecLog = ""
	if args.EnvJson == "" {
		args.EnvJson = "{}"
	}
	if len(args.TagIds) > 0 {
		if err := m.validateTagIdsExist(ctx, args.TagIds); err != nil {
			return 0, err
		}
	}
	var affect int64
	_, err := m.engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		s = s.Context(ctx)
		a, err := s.InsertOne(&args.DevProcess)
		if err != nil {
			return nil, err
		}
		affect = a
		err = replaceTagLinksSession(s, args.Id, dedupeTagIds(args.TagIds))
		return nil, err
	})
	return affect, err
}
