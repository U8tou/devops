package dev_process

import (
	devimpl "devops/impl/dev_process"
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func TagAdd(c *fiber.Ctx) error {
	var rq TagAddReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	if _, err := constx.GetLoginId(c); err != nil {
		return err
	}
	id, err := devimpl.Impl().TagAdd(c.Context(), rq.Name)
	if err != nil {
		return err
	}
	return r.Resp(c, TagAddResp{Id: datacv.IntToStr(id)})
}
