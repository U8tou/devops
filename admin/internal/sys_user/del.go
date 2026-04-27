package sysuser

import (
	"admin/internal/datapermctx"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_删除
// @Description	删除记录\n权限标识：sys:user:del
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			ids	body		DelReq						true	"IDs"
// @Success		200	{object}	r.MyResp{data=DelResp}	"Success"
// @Failure		500	{object}	r.MyResp					"Failure"
// @Router			/sys_user/del [delete]
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
	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	var iSysUserImpl = sysuser.Impl()
	if !isRoot {
		for _, idStr := range rq.Ids {
			row, gerr := iSysUserImpl.Get(c.Context(), idStr)
			if gerr != nil {
				return gerr
			}
			if row == nil || row.Id == 0 {
				continue
			}
			if row.UserType == model.UserTypeRoot {
				return errs.ERR_SYS_403
			}
			if err := datapermctx.CheckSubjectUserWith(c, loginId, isRoot, row.Id); err != nil {
				return err
			}
		}
	}
	// 处理
	affect, err := iSysUserImpl.Del(c.Context(), datacv.StrSliceToInt64Slice(rq.Ids))
	if err != nil {
		return err
	}
	// 响应
	var rp DelResp
	rp.Affect = datacv.IntToStr(affect)
	return r.Resp(c, rp)
}
