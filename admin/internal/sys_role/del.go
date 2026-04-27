package sysrole

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_删除
// @Description	删除角色\n权限标识：sys:role:del
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		DelReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=DelResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_role/del [delete]
func Del(c *fiber.Ctx) error {
	// 请求
	var rq DelReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var iSysRoleImpl = sysrole.Impl()
	affect, err := iSysRoleImpl.Del(c.Context(), datacv.StrSliceToInt64Slice(rq.Ids))
	if err != nil {
		return err
	}
	// 响应
	return r.Resp(c, DelResp{Affect: datacv.IntToStr(affect)})
}
