package syspost

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	syspost "system/impl/sys_post"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		岗位_编辑
// @Description	编辑记录\n权限标识：sys:post:edit
// @Tags			SysPostApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			message	body		EditReq	true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=EditResp}
// @Router			/sys_post/edit [put]
func Edit(c *fiber.Ctx) error {
	var rq EditReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}
	var args model.SysPostDto
	args.Id = datacv.StrToInt(rq.Id)
	args.Name = rq.Name
	args.Sort = rq.Sort
	args.Status = rq.Status
	args.UpdateBy = loginId
	impl := syspost.Impl()
	affect, err := impl.Edit(c.Context(), &args)
	if err != nil {
		return err
	}
	return r.Resp(c, EditResp{Affect: datacv.IntToStr(affect)})
}
