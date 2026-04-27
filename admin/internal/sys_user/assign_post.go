package sysuser

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_分配岗位
// @Description	为用户分配岗位\n权限标识：sys:user:assign_post
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AssignPostReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AssignPostResp}	"Success"
// @Failure		500		{object}	r.MyResp							"Failure"
// @Router			/sys_user/assign_post [post]
func AssignPost(c *fiber.Ctx) error {
	var rq AssignPostReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	iSysUserImpl := sysuser.Impl()
	affect, err := iSysUserImpl.AssignPost(c.Context(), datacv.StrToInt(rq.UserId), datacv.StrSliceToInt64Slice(rq.PostIds))
	if err != nil {
		return err
	}
	return r.Resp(c, AssignPostResp{Affect: datacv.IntToStr(affect)})
}
