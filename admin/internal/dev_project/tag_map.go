package dev_project

import (
	"devops/model"
	"pkg/tools/datacv"
)

func modelTagsToItems(tags []model.DevProjectTagBrief) []TagItem {
	if len(tags) == 0 {
		return nil
	}
	out := make([]TagItem, len(tags))
	for i, t := range tags {
		out[i] = TagItem{Id: datacv.IntToStr(t.Id), Name: t.Name}
	}
	return out
}
