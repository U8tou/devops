package sysauth

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	"strconv"
	sysauth "system/impl/sys_auth"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统认证_登入信息
// @Description	登入信息
// @Tags			SysAuthApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Success		200	{object}	r.MyResp{data=InfoResp}	"Success"
// @Failure		500	{object}	r.MyResp				"Failure"
// @Router			/sys_auth/info [get]
func Info(c *fiber.Ctx) error {
	// 请求
	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}
	// 处理
	sysAuthImpl := sysauth.Impl()
	loginInfo, err := sysAuthImpl.Info(c.Context(), datacv.IntToStr(loginId))
	if err != nil {
		return err
	}
	// 响应
	infoResp := InfoResp{
		UserId:     strconv.FormatInt(loginInfo.UserId, 10),
		UserName:   loginInfo.UserName,
		NickName:   loginInfo.NickName,
		UserType:   loginInfo.UserType,
		Email:      loginInfo.Email,
		PhoneArea:  loginInfo.PhoneArea,
		Phone:      loginInfo.Phone,
		Sex:        loginInfo.Sex,
		Avatar:     loginInfo.Avatar,
		Status:     loginInfo.Status,
		Address:    loginInfo.Address,
		Remark:     loginInfo.Remark,
		CreateTime: loginInfo.CreateTime,
		CreateBy:   loginInfo.CreateBy,
		UpdateTime: loginInfo.UpdateTime,
		UpdateBy:   loginInfo.UpdateBy,
		Depts:      loginInfo.Depts,
		Roles:      loginInfo.Roles,
		Posts:      loginInfo.Posts,
		Buttons:    loginInfo.Buttons,
		Menus:      loginInfo.Menus,
	}
	return r.Resp(c, infoResp)
}
