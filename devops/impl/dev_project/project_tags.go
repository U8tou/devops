package devproject

import (
	"context"
	"devops/model"

	"xorm.io/xorm"
)

func dedupeProjectTagIds(tagIds []int64) []int64 {
	if len(tagIds) <= 1 {
		return tagIds
	}
	seen := make(map[int64]struct{}, len(tagIds))
	out := make([]int64, 0, len(tagIds))
	for _, id := range tagIds {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func replaceProjectTagLinksSession(s *xorm.Session, projectId int64, tagIds []int64) error {
	if _, err := s.Where("project_id = ?", projectId).Delete(&model.DevProjectTagLink{}); err != nil {
		return err
	}
	for _, tid := range tagIds {
		if _, err := s.Insert(&model.DevProjectTagLink{ProjectId: projectId, TagId: tid}); err != nil {
			return err
		}
	}
	return nil
}

func (m *DevProjectImpl) validateProjectTagIdsExist(ctx context.Context, tagIds []int64) error {
	if len(tagIds) == 0 {
		return nil
	}
	cnt, err := m.engine.Context(ctx).Table(&model.DevProjectTag{}).In("id", tagIds).Count()
	if err != nil {
		return err
	}
	if cnt != int64(len(tagIds)) {
		return errTagIds
	}
	return nil
}

func collectProjectTagIds(links []model.DevProjectTagLink) []int64 {
	seen := make(map[int64]struct{})
	var ids []int64
	for _, l := range links {
		if _, ok := seen[l.TagId]; ok {
			continue
		}
		seen[l.TagId] = struct{}{}
		ids = append(ids, l.TagId)
	}
	return ids
}

func (m *DevProjectImpl) projectTagIdToNameMap(ctx context.Context, tagIds []int64) (map[int64]string, error) {
	out := make(map[int64]string)
	if len(tagIds) == 0 {
		return out, nil
	}
	var tags []model.DevProjectTag
	if err := m.engine.Context(ctx).In("id", tagIds).Find(&tags); err != nil {
		return nil, err
	}
	for _, t := range tags {
		out[t.Id] = t.Name
	}
	return out, nil
}

// LoadProjectTags 单项目标签（含孤儿：name 为空）
func (m *DevProjectImpl) LoadProjectTags(ctx context.Context, projectId int64) ([]model.DevProjectTagBrief, error) {
	var links []model.DevProjectTagLink
	err := m.engine.Context(ctx).Where("project_id = ?", projectId).Find(&links)
	if err != nil || len(links) == 0 {
		return nil, err
	}
	tagNames, err := m.projectTagIdToNameMap(ctx, collectProjectTagIds(links))
	if err != nil {
		return nil, err
	}
	out := make([]model.DevProjectTagBrief, 0, len(links))
	for _, l := range links {
		name := tagNames[l.TagId]
		out = append(out, model.DevProjectTagBrief{Id: l.TagId, Name: name})
	}
	return out, nil
}

// BatchLoadProjectTags 批量加载多项目标签
func (m *DevProjectImpl) BatchLoadProjectTags(ctx context.Context, projectIds []int64) (map[int64][]model.DevProjectTagBrief, error) {
	res := make(map[int64][]model.DevProjectTagBrief)
	if len(projectIds) == 0 {
		return res, nil
	}
	var links []model.DevProjectTagLink
	err := m.engine.Context(ctx).In("project_id", projectIds).Find(&links)
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return res, nil
	}
	tagNames, err := m.projectTagIdToNameMap(ctx, collectProjectTagIds(links))
	if err != nil {
		return nil, err
	}
	for _, l := range links {
		name := tagNames[l.TagId]
		res[l.ProjectId] = append(res[l.ProjectId], model.DevProjectTagBrief{Id: l.TagId, Name: name})
	}
	return res, nil
}

// DeleteLinksByProjectIds 删除项目时清理关联
func (m *DevProjectImpl) DeleteLinksByProjectIds(ctx context.Context, projectIds []int64) error {
	if len(projectIds) == 0 {
		return nil
	}
	_, err := m.engine.Context(ctx).In("project_id", projectIds).Delete(&model.DevProjectTagLink{})
	return err
}
