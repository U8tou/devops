package sysuser

import (
	"admin/internal/datapermctx"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"
	"system/model"
	"time"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_编辑
// @Description	编辑记录\n权限标识：sys:user:edit
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			message	body		EditReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=EditResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_user/edit [put]
func Edit(c *fiber.Ctx) error {
	// 请求
	var rq EditReq
	err := c.BodyParser(&rq)
	if err != nil {
		return err
	}
	err = validator.Struct(rq)
	if err != nil {
		return err
	}

	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	tp := time.Now().Unix()

	var iSysUserImpl = sysuser.Impl()
	existing, err := iSysUserImpl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if existing == nil || existing.Id == 0 {
		return errs.ERR_DB_NO_EXIST
	}
	if existing.UserType == model.UserTypeRoot && !isRoot {
		return errs.ERR_SYS_403
	}
	if err := datapermctx.CheckSubjectUserWith(c, loginId, isRoot, existing.Id); err != nil {
		return err
	}

	// 处理
	var args model.SysUserDto
	args.Id = datacv.StrToInt(rq.Id)
	args.UserName = rq.UserName
	args.NickName = rq.NickName
	args.UserType = model.UserTypeNormal
	if existing.UserType == model.UserTypeRoot {
		args.UserType = model.UserTypeRoot
	}
	args.Email = rq.Email
	args.PhoneArea = rq.PhoneArea
	args.Phone = rq.Phone
	args.Sex = rq.Sex
	args.Avatar = rq.Avatar
	args.Password = rq.Password
	args.Status = rq.Status
	args.Address = rq.Address
	args.Remark = rq.Remark
	args.UpdateTime = tp
	args.UpdateBy = loginId
	affect, err := iSysUserImpl.Edit(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp EditResp
	rp.Affect = datacv.IntToStr(affect)
	return r.Resp(c, rp)
}
