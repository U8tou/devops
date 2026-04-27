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

type SetCronEnabledReq struct {
	Id      string `json:"id" validate:"required,number"`
	Enabled bool   `json:"enabled"`
}

type SetCronEnabledResp struct {
	Affect string `json:"affect"`
}

func SetCronEnabled(c *fiber.Ctx) error {
	var rq SetCronEnabledReq
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
	var en int8
	if rq.Enabled {
		en = 1
	}
	id := datacv.StrToInt(rq.Id)
	affect, err := impl.SetCronEnabled(c.Context(), id, en, loginId)
	if err != nil {
		return err
	}
	SyncCronForProcess(id)
	return r.Resp(c, SetCronEnabledResp{Affect: datacv.IntToStr(affect)})
}
