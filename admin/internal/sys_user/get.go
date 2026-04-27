package sysuser

import (
	"admin/internal/datapermctx"
	"pkg/conf"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_查找
// @Description	获取详情\n权限标识：sys:user:get
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			id	query		GetReq						true	"ID"
// @Success		200	{object}	r.MyResp{data=GetResp}	"Success"
// @Failure		500	{object}	r.MyResp					"Failure"
// @Router			/sys_user/get [get]
func Get(c *fiber.Ctx) error {
	// 请求
	var rq GetReq
	rq.Id = c.Query("id")
	err := validator.Struct(rq)
	if err != nil {
		return errs.Args(err)
	}
	// 处理
	var iSysUserImpl = sysuser.Impl()
	row, err := iSysUserImpl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	if row.UserType == model.UserTypeRoot && !isRoot {
		return errs.ERR_SYS_403
	}
	if err := datapermctx.CheckSubjectUserWith(c, loginId, isRoot, row.Id); err != nil {
		return err
	}
	// 响应
	var rp GetResp
	rp.Id = datacv.IntToStr(row.Id)
	rp.UserName = row.UserName
	rp.NickName = row.NickName
	rp.UserType = row.UserType
	rp.Email = row.Email
	rp.PhoneArea = row.PhoneArea
	rp.Phone = row.Phone
	rp.Sex = row.Sex
	rp.Avatar = conf.FileUrl(row.Avatar)
	rp.Status = row.Status
	rp.Address = row.Address
	rp.Remark = row.Remark
	rp.CreateTime = datacv.IntToStr(row.CreateTime)
	rp.CreateBy = datacv.IntToStr(row.CreateBy)
	rp.UpdateTime = datacv.IntToStr(row.UpdateTime)
	rp.UpdateBy = datacv.IntToStr(row.UpdateBy)
	rp.Depts = row.Depts
	rp.Roles = row.Roles
	rp.Posts = row.Posts
	return r.Resp(c, rp)
}
