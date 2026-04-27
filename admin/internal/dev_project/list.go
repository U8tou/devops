package dev_project

import (
	"strconv"

	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_project"
	"devops/model"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func List(c *fiber.Ctx) error {
	var rq ListReq
	if err := c.QueryParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	var args model.DevProjectDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.FilterName = rq.Name
	if rq.Status != "" {
		s, err := strconv.ParseInt(rq.Status, 10, 8)
		if err != nil {
			return err
		}
		v := int8(s)
		args.FilterStatus = &v
	}
	args.TagFilterIds = parseCommaInt64s(rq.TagIds)
	args.TagFilterOther = parseTagOther(rq.TagOther)
	args.TagFilterExclude = parseTagModeExclude(rq.TagMode)
	if err := datapermctx.ApplyDevProjectDto(c, &args); err != nil {
		return err
	}
	impl := devimpl.Impl()
	total, rows, err := impl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	var rp ListResp
	rp.Total = datacv.IntToStr(total)
	rp.Rows = make([]ListBody, len(rows))
	for i, row := range rows {
		rp.Rows[i] = ListBody{
			Id:               datacv.IntToStr(row.Id),
			Name:             row.Name,
			Status:           strconv.Itoa(int(row.Status)),
			Progress:         strconv.Itoa(int(row.Progress)),
			VersionChangelog: row.VersionChangelog,
			CreateTime:       datacv.IntToStr(row.CreateTime),
			UpdateTime:       datacv.IntToStr(row.UpdateTime),
			Tags:             modelTagsToItems(row.Tags),
		}
	}
	return r.Resp(c, rp)
}
