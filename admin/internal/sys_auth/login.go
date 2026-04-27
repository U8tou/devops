package sysauth

import (
	r "pkg/resp"
	"pkg/validator"
	sysauth "system/impl/sys_auth"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_登入
// @Description	用户登陆
// @Tags			SysAuthApi
// @Accept			json
// @Produce		json
// @Param			args	body		LoginReq					true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=LoginResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_auth/login [post]
func Login(c *fiber.Ctx) error {
	// 请求
	var req LoginReq
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	if err := validator.Struct(req); err != nil {
		return err
	}
	// 处理
	sysAuthImpl := sysauth.Impl()
	tokenInfo, err := sysAuthImpl.Login(c.Context(), &model.LoginReq{
		UserName: req.UserName,
		Password: req.Password,
		CodeId:   req.CodeId,
		Code:     req.Code,
	})
	if err != nil {
		return err
	}

	// 响应
	resp := LoginResp{
		Token:        tokenInfo.Token,
		RefreshToken: tokenInfo.RefreshToken,
	}
	return r.Resp(c, resp)
}
