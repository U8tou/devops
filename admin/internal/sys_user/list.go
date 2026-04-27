package sysuser

import (
	"admin/internal/datapermctx"
	"pkg/conf"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysuser "system/impl/sys_user"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统用户_分页
// @Description	分页列表\n权限标识：sys:user:list
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		ListReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_user/list [get]
func List(c *fiber.Ctx) error {
	// 请求
	var rq ListReq
	if err := c.QueryParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var args model.SysUserDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.UserName = rq.UserName
	args.NickName = rq.NickName
	args.Email = rq.Email
	args.Phone = rq.Phone
	args.Sex = rq.Sex
	args.Status = rq.Status
	args.DeptId = datacv.StrToInt(rq.DeptId)
	if len(rq.TimeRange) == 2 {
		args.CreateTimeStart = datacv.StrToInt(rq.TimeRange[0])
		args.CreateTimeEnd = datacv.StrToInt(rq.TimeRange[1])
	}
	if err := datapermctx.ApplyUserDto(c, &args); err != nil {
		return err
	}
	var iSysUserImpl = sysuser.Impl()
	total, rows, err := iSysUserImpl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp ListResp
	rp.Total = datacv.IntToStr(total)
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
		body.Depts = row.Depts
		body.Roles = row.Roles
		body.Posts = row.Posts
		rp.Rows[idx] = body
	}
	return r.Resp(c, rp)
}
