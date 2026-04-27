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

// @Summary		岗位_新增
// @Description	新增记录\n权限标识：sys:post:add
// @Tags			SysPostApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AddReq	true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AddResp}
// @Router			/sys_post/add [post]
func Add(c *fiber.Ctx) error {
	var rq AddReq
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
	args.Name = rq.Name
	args.Sort = rq.Sort
	args.Status = rq.Status
	args.CreateBy = loginId
	args.UpdateBy = loginId
	impl := syspost.Impl()
	affect, err := impl.Add(c.Context(), &args)
	if err != nil {
		return err
	}
	return r.Resp(c, AddResp{Affect: datacv.IntToStr(affect)})
}
