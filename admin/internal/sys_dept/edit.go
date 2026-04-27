package sysdept

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysdept "system/impl/sys_dept"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统部门表_编辑
// @Description	编辑记录\n权限标识：sys:dept:edit
// @Tags			SysDeptApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			message	body		EditReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=EditResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_dept/edit [put]
func Edit(c *fiber.Ctx) error {
	// 请求
	var rq EditReq
	err := c.BodyParser(&rq)
	if err != nil {
		return err
	}
	err = validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var args model.SysDeptDto
	args.Id = datacv.StrToInt(rq.Id)
	args.Pid = datacv.StrToInt(rq.Pid)
	args.Name = rq.Name
	args.Profile = rq.Profile
	args.Leader = rq.Leader
	args.Phone = rq.Phone
	args.Email = rq.Email
	args.Sort = rq.Sort
	args.Status = rq.Status
	var iSysDeptImpl = sysdept.Impl()
	affect, err := iSysDeptImpl.Edit(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp EditResp
	rp.Affect = datacv.IntToStr(affect)
	return r.Resp(c, rp)
}
