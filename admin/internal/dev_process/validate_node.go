package dev_process

import (
	"pkg/devflow"
	r "pkg/resp"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ValidateNodeReq struct {
	Kind   string                 `json:"kind" validate:"required"`
	Params map[string]interface{} `json:"params"`
}

type ValidateNodeResp struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func ValidateNode(c *fiber.Ctx) error {
	var rq ValidateNodeReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	if rq.Params == nil {
		rq.Params = map[string]interface{}{}
	}
	ctx := c.UserContext()
	if ctx == nil {
		ctx = c.Context()
	}
	err := devflow.ValidateNode(ctx, rq.Kind, rq.Params)
	if err != nil {
		return r.Resp(c, ValidateNodeResp{
			Ok:      false,
			Message: err.Error(),
			Detail:  "",
		})
	}
	return r.Resp(c, ValidateNodeResp{
		Ok:      true,
		Message: "ok",
		Detail:  "",
	})
}
