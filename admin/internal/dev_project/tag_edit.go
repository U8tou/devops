package dev_project

import (
	devimpl "devops/impl/dev_project"
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func TagEdit(c *fiber.Ctx) error {
	var rq TagEditReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	if _, err := constx.GetLoginId(c); err != nil {
		return err
	}
	affect, err := devimpl.Impl().TagEdit(c.Context(), datacv.StrToInt(rq.Id), rq.Name)
	if err != nil {
		return err
	}
	return r.Resp(c, map[string]string{"affect": datacv.IntToStr(affect)})
}
