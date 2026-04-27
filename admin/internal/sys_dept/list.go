package sysdept

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysdept "system/impl/sys_dept"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统部门表_分页
// @Description	分页列表\n权限标识：sys:dept:list
// @Tags			SysDeptApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		ListReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_dept/list [get]
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
	var args model.SysDeptDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.Pid = datacv.StrToInt(rq.Pid)
	args.Name = rq.Name
	args.Profile = rq.Profile
	args.Leader = rq.Leader
	args.Phone = rq.Phone
	args.Email = rq.Email
	args.Sort = rq.Sort
	args.Status = rq.Status
	var iSysDeptImpl = sysdept.Impl()
	total, rows, err := iSysDeptImpl.List(c.Context(), &args)
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
		body.Pid = datacv.IntToStr(row.Pid)
		body.Name = row.Name
		body.Profile = row.Profile
		body.Leader = row.Leader
		body.Phone = row.Phone
		body.Email = row.Email
		body.Sort = row.Sort
		body.Status = row.Status
		body.CreateTime = datacv.IntToStr(row.CreateTime)
		body.UpdateTime = datacv.IntToStr(row.UpdateTime)
		rp.Rows[idx] = body
	}
	return r.Resp(c, rp)
}
