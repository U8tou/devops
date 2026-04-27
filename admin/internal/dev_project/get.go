package dev_project

import (
	"strconv"

	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_project"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	var rq GetReq
	rq.Id = c.Query("id")
	if err := validator.Struct(rq); err != nil {
		return errs.Args(err)
	}
	impl := devimpl.Impl()
	row, err := impl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	if err := datapermctx.CheckCreateByWith(c, loginId, isRoot, row.CreateBy); err != nil {
		return err
	}
	rp := GetResp{
		Id:               datacv.IntToStr(row.Id),
		Name:             row.Name,
		Status:           strconv.Itoa(int(row.Status)),
		Progress:         strconv.Itoa(int(row.Progress)),
		VersionChangelog: row.VersionChangelog,
		MindJson:         row.MindJson,
		CreateTime:       datacv.IntToStr(row.CreateTime),
		CreateBy:         datacv.IntToStr(row.CreateBy),
		UpdateTime:       datacv.IntToStr(row.UpdateTime),
		UpdateBy:         datacv.IntToStr(row.UpdateBy),
		Tags:             modelTagsToItems(row.Tags),
	}
	return r.Resp(c, rp)
}
