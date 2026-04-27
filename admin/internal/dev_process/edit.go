package dev_process

import (
	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_process"
	"devops/model"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func Edit(c *fiber.Ctx) error {
	var rq EditReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	impl := devimpl.Impl()
	existing, err := impl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if existing == nil {
		return errs.ERR_DB_NO_EXIST
	}
	if err := datapermctx.CheckCreateByWith(c, loginId, isRoot, existing.CreateBy); err != nil {
		return err
	}
	var args model.DevProcessDto
	args.Id = datacv.StrToInt(rq.Id)
	args.Code = rq.Code
	args.Remark = rq.Remark
	args.CronEnabled = rq.CronEnabled
	args.CronExpr = rq.CronExpr
	args.UpdateBy = loginId
	args.TagIds = rq.TagIds
	affect, err := impl.Edit(c.Context(), &args)
	if err != nil {
		return err
	}
	SyncCronForProcess(args.Id)
	return r.Resp(c, EditResp{Affect: datacv.IntToStr(affect)})
}
