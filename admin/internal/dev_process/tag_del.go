package dev_process

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	devimpl "devops/impl/dev_process"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func TagDel(c *fiber.Ctx) error {
	var rq TagDelReq
	rq.Id = c.Query("id")
	if err := validator.Struct(rq); err != nil {
		return err
	}
	if _, err := constx.GetLoginId(c); err != nil {
		return err
	}
	affect, err := devimpl.Impl().TagDel(c.Context(), datacv.StrToInt(rq.Id))
	if err != nil {
		return err
	}
	return r.Resp(c, map[string]string{"affect": datacv.IntToStr(affect)})
}
