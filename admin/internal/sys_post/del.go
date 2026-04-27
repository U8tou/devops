package syspost

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	syspost "system/impl/sys_post"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		岗位_删除
// @Description	删除记录\n权限标识：sys:post:del
// @Tags			SysPostApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			ids	body		DelReq	true	"IDs"
// @Success		200	{object}	r.MyResp{data=DelResp}
// @Router			/sys_post/del [delete]
func Del(c *fiber.Ctx) error {
	var rq DelReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	impl := syspost.Impl()
	affect, err := impl.Del(c.Context(), datacv.StrSliceToInt64Slice(rq.Ids))
	if err != nil {
		return err
	}
	return r.Resp(c, DelResp{Affect: datacv.IntToStr(affect)})
}
