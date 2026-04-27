package sysdept

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysdept "system/impl/sys_dept"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统部门表_删除
// @Description	删除记录\n权限标识：sys:dept:del
// @Tags			SysDeptApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			ids	query		DelReq						true	"IDs"
// @Success		200	{object}	r.MyResp{data=DelResp}	"Success"
// @Failure		500	{object}	r.MyResp					"Failure"
// @Router			/sys_dept/del [delete]
func Del(c *fiber.Ctx) error {
	// 请求
	var rq DelReq
	err := c.BodyParser(&rq)
	if err != nil {
		return err
	}
	err = validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var iSysDeptImpl = sysdept.Impl()
	affect, err := iSysDeptImpl.Del(c.Context(), datacv.StrSliceToInt64Slice(rq.Ids))
	if err != nil {
		return err
	}
	// 响应
	var rp DelResp
	rp.Affect = datacv.IntToStr(affect)
	return r.Resp(c, rp)
}
