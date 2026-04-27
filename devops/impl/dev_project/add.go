package devproject

import (
	"context"
	"devops/model"
	"strings"

	"github.com/yitter/idgenerator-go/idgen"
	"xorm.io/xorm"
)

func clampProgress(p int8) int8 {
	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}

func (m *DevProjectImpl) Add(ctx context.Context, args *model.DevProjectDto) (int64, error) {
	v := &args.DevProject
	v.Id = idgen.NextId()
	v.DeleteTime = 0
	v.Name = strings.TrimSpace(v.Name)
	v.Progress = clampProgress(v.Progress)
	if v.MindJson == "" {
		v.MindJson = "{}"
	}
	if len(args.TagIds) > 0 {
		if err := m.validateProjectTagIdsExist(ctx, args.TagIds); err != nil {
			return 0, err
		}
	}
	var affect int64
	_, err := m.engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		s = s.Context(ctx)
		a, err := s.InsertOne(v)
		if err != nil {
			return nil, err
		}
		affect = a
		err = replaceProjectTagLinksSession(s, v.Id, dedupeProjectTagIds(args.TagIds))
		return nil, err
	})
	return affect, err
}
