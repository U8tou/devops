package sysuser

import (
	"pkg/constx"
	r "pkg/resp"
	sysuser "system/impl/sys_user"
	"system/model"
	"time"

	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_新增
// @Description	新增记录\n权限标识：sys:user:add
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AddReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AddResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_user/add [post]
func Add(c *fiber.Ctx) error {
	// 请求
	var rq AddReq
	err := c.BodyParser(&rq)
	if err != nil {
		return err
	}
	err = validator.Struct(rq)
	if err != nil {
		return err
	}

	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}
	tp := time.Now().Unix()

	// 处理
	var args model.SysUserDto
	args.UserName = rq.UserName
	args.NickName = rq.NickName
	args.UserType = model.UserTypeNormal
	args.Email = rq.Email
	args.PhoneArea = rq.PhoneArea
	args.Phone = rq.Phone
	args.Sex = rq.Sex
	args.Avatar = rq.Avatar
	args.Password = rq.Password
	args.Status = rq.Status
	args.Address = rq.Address
	args.Remark = rq.Remark
	args.CreateTime = tp
	args.CreateBy = loginId
	args.UpdateTime = tp
	args.UpdateBy = loginId
	var iSysUserImpl = sysuser.Impl()
	affect, err := iSysUserImpl.Add(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp AddResp
	rp.Affect = datacv.IntToStr(affect)
	return r.Resp(c, rp)
}
