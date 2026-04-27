package sysrole

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_分配菜单权限
// @Description	为角色分配菜单权限\n权限标识：sys:role:assign_menu
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AssignMenuReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AssignMenuResp}	"Success"
// @Failure		500		{object}	r.MyResp							"Failure"
// @Router			/sys_role/assign_menu [post]
func AssignMenu(c *fiber.Ctx) error {
	// 请求
	var rq AssignMenuReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var iSysRoleImpl = sysrole.Impl()
	affect, err := iSysRoleImpl.AssignMenu(c.Context(), datacv.StrToInt(rq.RoleId), datacv.StrSliceToInt64Slice(rq.MenuIds), rq.MenuLinkage)
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, AssignMenuResp{Affect: datacv.IntToStr(affect)})
}
