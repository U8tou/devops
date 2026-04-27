package sysdept

import (
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysdept "system/impl/sys_dept"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统部门表_查找
// @Description	获取详情\n权限标识：sys:dept:get
// @Tags			SysDeptApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			id	query		GetReq						true	"ID"
// @Success		200	{object}	r.MyResp{data=GetResp}	"Success"
// @Failure		500	{object}	r.MyResp					"Failure"
// @Router			/sys_dept/get [get]
func Get(c *fiber.Ctx) error {
	// 请求
	var rq GetReq
	rq.Id = c.Query("id")
	err := validator.Struct(rq)
	if err != nil {
		return errs.Args(err)
	}
	// 处理
	var iSysDeptImpl = sysdept.Impl()
	row, err := iSysDeptImpl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	// 响应
	var rp GetResp
	rp.Id = datacv.IntToStr(row.Id)
	rp.Pid = datacv.IntToStr(row.Pid)
	rp.Name = row.Name
	rp.Profile = row.Profile
	rp.Leader = row.Leader
	rp.Phone = row.Phone
	rp.Email = row.Email
	rp.Sort = row.Sort
	rp.Status = row.Status
	rp.CreateTime = datacv.IntToStr(row.CreateTime)
	rp.CreateBy = datacv.IntToStr(row.CreateBy)
	rp.UpdateTime = datacv.IntToStr(row.UpdateTime)
	rp.UpdateBy = datacv.IntToStr(row.UpdateBy)
	return r.Resp(c, rp)
}
