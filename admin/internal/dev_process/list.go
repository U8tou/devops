package dev_process

import (
	"strconv"

	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_process"
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
	var args model.DevProcessDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.Code = rq.Code
	args.Remark = rq.Remark
	args.TagFilterIds = parseCommaInt64s(rq.TagIds)
	args.TagFilterOther = parseTagOther(rq.TagOther)
	args.TagFilterExclude = parseTagModeExclude(rq.TagMode)
	if err := datapermctx.ApplyDevProcessDto(c, &args); err != nil {
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
			Id:                 datacv.IntToStr(row.Id),
			Code:               row.Code,
			Remark:             row.Remark,
			CronExpr:           row.CronExpr,
			CronEnabled:        strconv.Itoa(int(row.CronEnabled)),
			LastExecTime:       datacv.IntToStr(row.LastExecTime),
			LastExecDurationMs: datacv.IntToStr(row.LastExecDurationMs),
			LastExecResult:     normalizeListExecStatus(row.LastExecResult),
			CreateTime:         datacv.IntToStr(row.CreateTime),
			UpdateTime:         datacv.IntToStr(row.UpdateTime),
			Tags:               modelTagsToItems(row.Tags),
		}
	}
	return r.Resp(c, rp)
}
