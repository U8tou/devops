package sysuser

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_分配角色
// @Description	为用户分配角色\n权限标识：sys:user:assign_role
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AssignRoleReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AssignRoleResp}	"Success"
// @Failure		500		{object}	r.MyResp							"Failure"
// @Router			/sys_user/assign_role [post]
func AssignRole(c *fiber.Ctx) error {
	// 请求
	var rq AssignRoleReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var iSysUserImpl = sysuser.Impl()
	affect, err := iSysUserImpl.AssignRole(c.Context(), datacv.StrToInt(rq.UserId), datacv.StrSliceToInt64Slice(rq.RoleIds))
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, AssignRoleResp{Affect: datacv.IntToStr(affect)})
}
