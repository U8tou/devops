package sysauth

import (
	r "pkg/resp"
	sysauth "system/impl/sys_auth"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_Token刷新
// @Description	使用 RefreshToken 换取新的 AccessToken
// @Tags			SysAuthApi
// @Accept			json
// @Produce		json
// @Param			args	body		RefreshTokenReq				true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=LoginResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_auth/refresh_token [post]
func RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenReq
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	device := c.Get("X-Device", "web")
	sysAuthImpl := sysauth.Impl()
	tokenInfo, err := sysAuthImpl.RefreshToken(c.Context(), device, req.RefreshToken)
	if err != nil {
		return err
	}
	resp := LoginResp{
		Token:        tokenInfo.Token,
		RefreshToken: tokenInfo.RefreshToken,
	}
	return r.Resp(c, resp)
}
