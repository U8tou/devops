package sysauth

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysauth "system/impl/sys_auth"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_修改密码
// @Description	当前登录用户修改密码
// @Tags			SysAuthApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		ChangePasswordReq	true	"Request Parameters"
// @Success		200		{object}	r.MyResp	"Success"
// @Failure		500		{object}	r.MyResp	"Failure"
// @Router			/sys_auth/change_password [post]
func ChangePassword(c *fiber.Ctx) error {
	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}

	var req ChangePasswordReq
	if err = c.BodyParser(&req); err != nil {
		return err
	}

	sysAuthImpl := sysauth.Impl()
	if err = sysAuthImpl.ChangePassword(c.Context(), datacv.IntToStr(loginId), req.OldPassword, req.NewPassword); err != nil {
		return err
	}

	return r.Ok(c)
}
