package sysauth

import (
	"pkg/conf"
	"pkg/constx"
	r "pkg/resp"
	sysauth "system/impl/sys_auth"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_登出
// @Description	用户登出
// @Tags			SysAuthApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Router			/sys_auth/logout [post]
func Logout(c *fiber.Ctx) error {
	token := c.Get(conf.Auth.TokenName)
	if token == "" {
		return fiber.ErrUnauthorized
	}
	device := c.Get("X-Device", "web")
	_, _ = constx.GetLoginId(c)
	impl := sysauth.Impl()
	if err := impl.Logout(c.Context(), device, token); err != nil {
		return err
	}
	return r.Ok(c)
}
