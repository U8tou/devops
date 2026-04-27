package devproject

import (
	"context"
	"devops/model"
	"strings"
	"time"

	"pkg/errs"

	"xorm.io/xorm"
)

func (m *DevProjectImpl) Edit(ctx context.Context, args *model.DevProjectDto) (int64, error) {
	v := args.DevProject
	v.Name = strings.TrimSpace(v.Name)
	v.Progress = clampProgress(v.Progress)
	var old model.DevProject
	has, err := m.engine.Context(ctx).ID(v.Id).Get(&old)
	if err != nil {
		return 0, err
	}
	if !has {
		return 0, errs.ERR_DB_NO_EXIST
	}
	v.UpdateTime = time.Now().Unix()
	if len(args.TagIds) > 0 {
		if err := m.validateProjectTagIdsExist(ctx, args.TagIds); err != nil {
			return 0, err
		}
	}
	var affect int64
	_, err = m.engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		s = s.Context(ctx)
		a, err := s.ID(v.Id).Cols("name", "status", "progress", "version_changelog", "update_by", "update_time").Update(&v)
		if err != nil {
			return nil, err
		}
		affect = a
		err = replaceProjectTagLinksSession(s, v.Id, dedupeProjectTagIds(args.TagIds))
		return nil, err
	})
	return affect, err
}
