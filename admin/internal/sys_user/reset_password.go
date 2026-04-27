package sysuser

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_重置密码
// @Description	为指定用户重置密码\n权限标识：sys:user:reset_pwd
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		ResetPasswordReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=ResetPasswordResp}	"Success"
// @Failure		500		{object}	r.MyResp								"Failure"
// @Router			/sys_user/reset_pwd [post]
func ResetPassword(c *fiber.Ctx) error {
	var rq ResetPasswordReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}

	iSysUserImpl := sysuser.Impl()
	affect, err := iSysUserImpl.ResetPassword(c.Context(), datacv.StrToInt(rq.UserId), rq.Password)
	if err != nil {
		return err
	}
	return r.Resp(c, ResetPasswordResp{Affect: datacv.IntToStr(affect)})
}

