package dev_project

import (
	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_project"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func EditMind(c *fiber.Ctx) error {
	var rq EditMindReq
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
	id := datacv.StrToInt(rq.Id)
	affect, err := impl.EditMind(c.Context(), id, rq.MindJson, loginId)
	if err != nil {
		return err
	}
	return r.Resp(c, EditMindResp{Affect: datacv.IntToStr(affect)})
}
