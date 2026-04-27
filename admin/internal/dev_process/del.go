package dev_process

import (
	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_process"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type DelReq struct {
	Ids []string `json:"ids" validate:"required,dive,number"`
}

type DelResp struct {
	Affect string `json:"affect"`
}

func Del(c *fiber.Ctx) error {
	var rq DelReq
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
	if !isRoot {
		for _, idStr := range rq.Ids {
			row, gerr := impl.Get(c.Context(), idStr)
			if gerr != nil {
				return gerr
			}
			if row == nil {
				continue
			}
			if err := datapermctx.CheckCreateByWith(c, loginId, isRoot, row.CreateBy); err != nil {
				return err
			}
		}
	}
	ids := datacv.StrSliceToInt64Slice(rq.Ids)
	affect, err := impl.Del(c.Context(), ids)
	if err != nil {
		return err
	}
	for _, id := range ids {
		RemoveCronTask(id)
	}
	return r.Resp(c, DelResp{Affect: datacv.IntToStr(affect)})
}
