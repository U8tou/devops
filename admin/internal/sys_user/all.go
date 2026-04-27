package sysuser

import (
	"admin/internal/datapermctx"
	"pkg/conf"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_全部
// @Description	全部数据\n权限标识：sys:user:get
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Success		200	{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500	{object}	r.MyResp					"Failure"
// @Router			/sys_user/all [get]
func All(c *fiber.Ctx) error {
	// 处理
	var args model.SysUserDto
	args.Status = 1
	if err := datapermctx.ApplyUserDto(c, &args); err != nil {
		return err
	}
	var iSysUserImpl = sysuser.Impl()
	_, rows, err := iSysUserImpl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp ListResp
	rp.Total = datacv.IntToStr(int64(len(rows)))
	rp.Rows = make([]ListBody, len(rows))
	for idx, row := range rows {
		var body ListBody
		body.Id = datacv.IntToStr(row.Id)
		body.UserName = row.UserName
		body.NickName = row.NickName
		body.UserType = row.UserType
		body.Email = row.Email
		body.PhoneArea = row.PhoneArea
		body.Phone = row.Phone
		body.Sex = row.Sex
		body.Avatar = conf.FileUrl(row.Avatar)
		body.Status = row.Status
		body.Address = row.Address
		body.Remark = row.Remark
		body.CreateTime = datacv.IntToStr(row.CreateTime)
		body.CreateBy = datacv.IntToStr(row.CreateBy)
		body.UpdateTime = datacv.IntToStr(row.UpdateTime)
		body.UpdateBy = datacv.IntToStr(row.UpdateBy)
		rp.Rows[idx] = body
	}
	return r.Resp(c, rp)
}
