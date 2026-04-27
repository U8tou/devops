package sysrole

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_编辑
// @Description	编辑角色\n权限标识：sys:role:edit
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		EditReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=EditResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_role/edit [put]
func Edit(c *fiber.Ctx) error {
	// 请求
	var rq EditReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 获取当前用户ID
	loginId, _ := constx.GetLoginId(c)
	// 处理
	var args model.SysRoleDto
	args.Id = datacv.StrToInt(rq.Id)
	args.Name = rq.Name
	args.Role = rq.Role
	args.Status = rq.Status
	args.Sort = rq.Sort
	args.Remark = rq.Remark
	args.UpdateBy = loginId
	var iSysRoleImpl = sysrole.Impl()
	affect, err := iSysRoleImpl.Edit(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, EditResp{Affect: datacv.IntToStr(affect)})
}
