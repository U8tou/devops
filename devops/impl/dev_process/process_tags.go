package devprocess

import (
	"context"
	"devops/model"

	"xorm.io/xorm"
)

func dedupeTagIds(tagIds []int64) []int64 {
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

func replaceTagLinksSession(s *xorm.Session, processId int64, tagIds []int64) error {
	if _, err := s.Where("process_id = ?", processId).Delete(&model.DevProcessTagLink{}); err != nil {
		return err
	}
	for _, tid := range tagIds {
		if _, err := s.Insert(&model.DevProcessTagLink{ProcessId: processId, TagId: tid}); err != nil {
			return err
		}
	}
	return nil
}

// ReplaceProcessTags 用「当前字典中的标签」覆盖用户所选；**保留**仍挂在流程上但字典已删的孤儿 tag_id
func (m *DevProcessImpl) ReplaceProcessTags(ctx context.Context, processId int64, tagIds []int64) error {
	tagIds = dedupeTagIds(tagIds)
	if len(tagIds) > 0 {
		if err := m.validateTagIdsExist(ctx, tagIds); err != nil {
			return err
		}
	}
	orphans, err := m.orphanTagIdsForProcess(ctx, processId)
	if err != nil {
		return err
	}
	merged := append(append([]int64{}, orphans...), tagIds...)
	merged = dedupeTagIds(merged)
	_, err = m.engine.Transaction(func(s *xorm.Session) (interface{}, error) {
		s = s.Context(ctx)
		err := replaceTagLinksSession(s, processId, merged)
		return nil, err
	})
	return err
}

func (m *DevProcessImpl) orphanTagIdsForProcess(ctx context.Context, processId int64) ([]int64, error) {
	var links []model.DevProcessTagLink
	if err := m.engine.Context(ctx).Where("process_id = ?", processId).Find(&links); err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, nil
	}
	ids := collectTagIds(links)
	var found []model.DevProcessTag
	if err := m.engine.Context(ctx).In("id", ids).Cols("id").Find(&found); err != nil {
		return nil, err
	}
	foundSet := make(map[int64]struct{}, len(found))
	for _, t := range found {
		foundSet[t.Id] = struct{}{}
	}
	var orphans []int64
	for _, id := range ids {
		if _, ok := foundSet[id]; !ok {
			orphans = append(orphans, id)
		}
	}
	return orphans, nil
}

func (m *DevProcessImpl) validateTagIdsExist(ctx context.Context, tagIds []int64) error {
	if len(tagIds) == 0 {
		return nil
	}
	cnt, err := m.engine.Context(ctx).Table(&model.DevProcessTag{}).In("id", tagIds).Count()
	if err != nil {
		return err
	}
	if cnt != int64(len(tagIds)) {
		return errTagIds
	}
	return nil
}

// LoadProcessTags 单流程标签（含孤儿：name 为空）
func (m *DevProcessImpl) LoadProcessTags(ctx context.Context, processId int64) ([]model.DevProcessTagBrief, error) {
	var links []model.DevProcessTagLink
	err := m.engine.Context(ctx).Where("process_id = ?", processId).Find(&links)
	if err != nil || len(links) == 0 {
		return nil, err
	}
	return m.linksToBriefs(ctx, links)
}

// BatchLoadProcessTags 批量加载多流程标签
func (m *DevProcessImpl) BatchLoadProcessTags(ctx context.Context, processIds []int64) (map[int64][]model.DevProcessTagBrief, error) {
	res := make(map[int64][]model.DevProcessTagBrief)
	if len(processIds) == 0 {
		return res, nil
	}
	var links []model.DevProcessTagLink
	err := m.engine.Context(ctx).In("process_id", processIds).Find(&links)
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return res, nil
	}
	tagNames, err := m.tagIdToNameMap(ctx, collectTagIds(links))
	if err != nil {
		return nil, err
	}
	for _, l := range links {
		name := tagNames[l.TagId]
		res[l.ProcessId] = append(res[l.ProcessId], model.DevProcessTagBrief{Id: l.TagId, Name: name})
	}
	return res, nil
}

func collectTagIds(links []model.DevProcessTagLink) []int64 {
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

func (m *DevProcessImpl) tagIdToNameMap(ctx context.Context, tagIds []int64) (map[int64]string, error) {
	out := make(map[int64]string)
	if len(tagIds) == 0 {
		return out, nil
	}
	var tags []model.DevProcessTag
	if err := m.engine.Context(ctx).In("id", tagIds).Find(&tags); err != nil {
		return nil, err
	}
	for _, t := range tags {
		out[t.Id] = t.Name
	}
	return out, nil
}

func (m *DevProcessImpl) linksToBriefs(ctx context.Context, links []model.DevProcessTagLink) ([]model.DevProcessTagBrief, error) {
	tagNames, err := m.tagIdToNameMap(ctx, collectTagIds(links))
	if err != nil {
		return nil, err
	}
	out := make([]model.DevProcessTagBrief, 0, len(links))
	for _, l := range links {
		name := tagNames[l.TagId]
		out = append(out, model.DevProcessTagBrief{Id: l.TagId, Name: name})
	}
	return out, nil
}

// DeleteLinksByProcessIds 删除流程时清理关联
func (m *DevProcessImpl) DeleteLinksByProcessIds(ctx context.Context, processIds []int64) error {
	if len(processIds) == 0 {
		return nil
	}
	_, err := m.engine.Context(ctx).In("process_id", processIds).Delete(&model.DevProcessTagLink{})
	return err
}
