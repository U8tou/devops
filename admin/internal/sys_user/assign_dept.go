package sysuser

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_分配部门
// @Description	为用户分配部门\n权限标识：sys:user:assign_dept
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AssignDeptReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AssignDeptResp}	"Success"
// @Failure		500		{object}	r.MyResp							"Failure"
// @Router			/sys_user/assign_dept [post]
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
	var iSysUserImpl = sysuser.Impl()
	affect, err := iSysUserImpl.AssignDept(c.Context(), datacv.StrToInt(rq.UserId), datacv.StrSliceToInt64Slice(rq.DeptIds))
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, AssignDeptResp{Affect: datacv.IntToStr(affect)})
}
