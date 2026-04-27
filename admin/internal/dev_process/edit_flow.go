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

func EditFlow(c *fiber.Ctx) error {
	var rq EditFlowReq
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
	args.Flow = rq.Flow
	args.UpdateBy = loginId
	affect, err := impl.EditFlow(c.Context(), &args)
	if err != nil {
		return err
	}
	return r.Resp(c, EditFlowResp{Affect: datacv.IntToStr(affect)})
}
