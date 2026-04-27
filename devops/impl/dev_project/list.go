package devproject

import (
	"context"
	"devops/model"
	"fmt"
	"strings"
)

func int64InPlaceholders(ids []int64) (string, []interface{}) {
	args := make([]interface{}, len(ids))
	parts := make([]string, len(ids))
	for i, id := range ids {
		parts[i] = "?"
		args[i] = id
	}
	return strings.Join(parts, ","), args
}

func (m *DevProjectImpl) List(ctx context.Context, args *model.DevProjectDto) (int64, []model.DevProjectVo, error) {
	query := m.engine.Context(ctx).Table(&model.DevProject{}).Omit("mind_json").Where("1 = 1")
	if args.DataScopeActive {
		if len(args.DataScopeDeptIds) == 0 {
			query.And("1 = 0")
		} else {
			ph, bind := int64InPlaceholders(args.DataScopeDeptIds)
			query.And(fmt.Sprintf("EXISTS (SELECT 1 FROM sys_user_dept u WHERE u.user_id = dev_project.create_by AND u.dept_id IN (%s))", ph), bind...)
		}
	}
	if args.FilterName != "" {
		query.And("name LIKE ?", "%"+args.FilterName+"%")
	}
	if args.FilterStatus != nil {
		query.And("status = ?", *args.FilterStatus)
	}

	hasTagIds := len(args.TagFilterIds) > 0
	if hasTagIds || args.TagFilterOther {
		var inner string
		var bind []interface{}
		if hasTagIds && args.TagFilterOther {
			ph, tagArgs := int64InPlaceholders(args.TagFilterIds)
			inner = fmt.Sprintf(
				"(id IN (SELECT DISTINCT project_id FROM dev_project_tag_link WHERE tag_id IN (%s)) OR EXISTS (SELECT 1 FROM dev_project_tag_link l WHERE l.project_id = dev_project.id AND l.tag_id NOT IN (SELECT id FROM dev_project_tag)))",
				ph,
			)
			bind = tagArgs
		} else if hasTagIds {
			ph, tagArgs := int64InPlaceholders(args.TagFilterIds)
			inner = fmt.Sprintf(
				"id IN (SELECT DISTINCT project_id FROM dev_project_tag_link WHERE tag_id IN (%s))",
				ph,
			)
			bind = tagArgs
		} else {
			inner = "EXISTS (SELECT 1 FROM dev_project_tag_link l WHERE l.project_id = dev_project.id AND l.tag_id NOT IN (SELECT id FROM dev_project_tag))"
		}
		if args.TagFilterExclude {
			if len(bind) > 0 {
				query.And("NOT ("+inner+")", bind...)
			} else {
				query.And("NOT (" + inner + ")")
			}
		} else {
			if len(bind) > 0 {
				query.And(inner, bind...)
			} else {
				query.And(inner)
			}
		}
	}

	var total int64
	var err error
	rows := make([]model.DevProjectVo, 0)
	if args.Size > 0 {
		total, err = query.Limit(args.Size, (args.Current-1)*args.Size).Desc("id").FindAndCount(&rows)
	} else {
		err = query.Desc("id").Find(&rows)
		if err == nil {
			total = int64(len(rows))
		}
	}
	if err != nil {
		return 0, nil, err
	}
	if len(rows) == 0 {
		return total, rows, nil
	}
	ids := make([]int64, len(rows))
	for i := range rows {
		ids[i] = rows[i].Id
	}
	tagMap, err := m.BatchLoadProjectTags(ctx, ids)
	if err != nil {
		return 0, nil, err
	}
	for i := range rows {
		rows[i].Tags = tagMap[rows[i].Id]
	}
	return total, rows, nil
}
