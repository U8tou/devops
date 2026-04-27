package sysrole

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_分配部门权限
// @Description	为角色分配部门权限\n权限标识：sys:role:assign_dept
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AssignDeptReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AssignDeptResp}	"Success"
// @Failure		500		{object}	r.MyResp							"Failure"
// @Router			/sys_role/assign_dept [post]
func AssignDept(c *fiber.Ctx) error {
	// 请求
	var rq AssignDeptReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var iSysRoleImpl = sysrole.Impl()
	affect, err := iSysRoleImpl.AssignDept(c.Context(), datacv.StrToInt(rq.RoleId), datacv.StrSliceToInt64Slice(rq.DeptIds), rq.DeptLinkage)
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, AssignDeptResp{Affect: datacv.IntToStr(affect)})
}
