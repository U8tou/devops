package dev_project

import (
	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_project"
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
	var args model.DevProjectDto
	args.Id = datacv.StrToInt(rq.Id)
	args.Name = rq.Name
	args.Status = rq.Status
	args.Progress = rq.Progress
	args.VersionChangelog = rq.VersionChangelog
	args.TagIds = rq.TagIds
	args.UpdateBy = loginId
	affect, err := impl.Edit(c.Context(), &args)
	if err != nil {
		return err
	}
	return r.Resp(c, EditResp{Affect: datacv.IntToStr(affect)})
}
