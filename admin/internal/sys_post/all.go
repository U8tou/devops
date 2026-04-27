package syspost

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	syspost "system/impl/sys_post"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// @Summary		岗位_全部
// @Description	全部数据（如下拉）\n权限标识：无需权限（仅需登录）
// @Tags			SysPostApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Success		200	{object}	r.MyResp{data=ListResp}
// @Router			/sys_post/all [get]
func All(c *fiber.Ctx) error {
	var args model.SysPostDto
	impl := syspost.Impl()
	_, rows, err := impl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	var rp ListResp
	rp.Total = datacv.IntToStr(int64(len(rows)))
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
