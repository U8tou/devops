package dev_process

import (
	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_process"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type EditEnvReq struct {
	Id  string            `json:"id" validate:"required,number"`
	Env map[string]string `json:"env"`
}

type EditEnvResp struct {
	Affect string `json:"affect"`
}

func EditEnv(c *fiber.Ctx) error {
	var rq EditEnvReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	if rq.Env == nil {
		rq.Env = map[string]string{}
	}
	raw, err := mapToEnvJSON(rq.Env)
	if err != nil {
		return errs.New(err.Error())
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
	affect, err := impl.EditEnv(c.Context(), datacv.StrToInt(rq.Id), raw, loginId)
	if err != nil {
		return err
	}
	return r.Resp(c, EditEnvResp{Affect: datacv.IntToStr(affect)})
}
