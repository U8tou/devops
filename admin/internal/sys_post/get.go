package syspost

import (
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"
	syspost "system/impl/sys_post"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		岗位_详情
// @Description	获取详情\n权限标识：sys:post:get
// @Tags			SysPostApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			id	query		GetReq	true	"ID"
// @Success		200	{object}	r.MyResp{data=GetResp}
// @Router			/sys_post/get [get]
func Get(c *fiber.Ctx) error {
	var rq GetReq
	rq.Id = c.Query("id")
	if err := validator.Struct(rq); err != nil {
		return errs.Args(err)
	}
	impl := syspost.Impl()
	row, err := impl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	rp := GetResp{
		Id:         datacv.IntToStr(row.Id),
		Name:       row.Name,
		Sort:       row.Sort,
		Status:     row.Status,
		CreateTime: datacv.IntToStr(row.CreateTime),
		CreateBy:   datacv.IntToStr(row.CreateBy),
		UpdateTime: datacv.IntToStr(row.UpdateTime),
		UpdateBy:   datacv.IntToStr(row.UpdateBy),
	}
	return r.Resp(c, rp)
}
