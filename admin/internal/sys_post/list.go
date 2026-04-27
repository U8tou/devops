package syspost

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	syspost "system/impl/sys_post"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		岗位_分页
// @Description	分页列表\n权限标识：sys:post:list
// @Tags			SysPostApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		ListReq	true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=ListResp}
// @Router			/sys_post/list [get]
func List(c *fiber.Ctx) error {
	var rq ListReq
	if err := c.QueryParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	var args model.SysPostDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.Name = rq.Name
	args.Sort = rq.Sort
	args.Status = rq.Status
	impl := syspost.Impl()
	total, rows, err := impl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	var rp ListResp
	rp.Total = datacv.IntToStr(total)
	rp.Rows = make([]ListBody, len(rows))
	for i, row := range rows {
		rp.Rows[i] = ListBody{
			Id:         datacv.IntToStr(row.Id),
			Name:       row.Name,
			Sort:       row.Sort,
			Status:     row.Status,
			CreateTime: datacv.IntToStr(row.CreateTime),
			UpdateTime: datacv.IntToStr(row.UpdateTime),
		}
	}
	return r.Resp(c, rp)
}
