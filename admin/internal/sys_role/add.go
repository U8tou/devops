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

// @Summary		系统角色_新增
// @Description	新增角色\n权限标识：sys:role:add
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AddReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AddResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_role/add [post]
func Add(c *fiber.Ctx) error {
	// 请求
	var rq AddReq
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
	args.Name = rq.Name
	args.Role = rq.Role
	args.Status = rq.Status
	args.Sort = rq.Sort
	args.Remark = rq.Remark
	args.CreateBy = loginId
	args.UpdateBy = loginId
	var iSysRoleImpl = sysrole.Impl()
	id, err := iSysRoleImpl.Add(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, AddResp{Id: datacv.IntToStr(id)})
}
