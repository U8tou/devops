package sysauth

import (
	r "pkg/resp"
	"pkg/validator"
	sysauth "system/impl/sys_auth"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_注册
// @Description	用户注册（公开接口，无需登录）
// @Tags			SysAuthApi
// @Accept			json
// @Produce		json
// @Param			args	body		RegisterReq					true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=RegisterResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_auth/register [post]
func Register(c *fiber.Ctx) error {
	var req RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := validator.Struct(req); err != nil {
		return err
	}

	impl := sysauth.Impl()
	if err := impl.Register(c.Context(), req.UserName, req.Password); err != nil {
		return err
	}

	return r.Resp(c, RegisterResp{UserName: req.UserName})
}
