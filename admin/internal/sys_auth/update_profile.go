package sysauth

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysauth "system/impl/sys_auth"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_更新个人信息
// @Description	当前登录用户更新昵称、邮箱、电话、性别、地址等
// @Tags			SysAuthApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		UpdateProfileReq	true	"Request Parameters"
// @Success		200		{object}	r.MyResp	"Success"
// @Failure		500		{object}	r.MyResp	"Failure"
// @Router			/sys_auth/profile [post]
func UpdateProfile(c *fiber.Ctx) error {
	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}

	var req UpdateProfileReq
	if err = c.BodyParser(&req); err != nil {
		return err
	}

	sysAuthImpl := sysauth.Impl()
	if err = sysAuthImpl.UpdateProfile(c.Context(), datacv.IntToStr(loginId), req.NickName, req.Email, req.PhoneArea, req.Phone, req.Sex, req.Address); err != nil {
		return err
	}

	return r.Resp(c, nil)
}
