package sysdept

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysdept "system/impl/sys_dept"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

/**
Notes: SysDept API
Time:  2025-04-29 10:48:59
*/

// @Summary		系统部门表_全部
// @Description	全部数据\n权限标识：无需权限（仅需登录）
// @Tags			SysDeptApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Success		200	{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500	{object}	r.MyResp					"Failure"
// @Router			/sys_dept/all [get]
func All(c *fiber.Ctx) error {
	// 处理
	var args model.SysDeptDto
	// args.Status = 1
	var iSysDeptImpl = sysdept.Impl()
	_, rows, err := iSysDeptImpl.List(c.Context(), &args)
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
