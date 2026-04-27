package sysmenu

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysmenu "system/impl/sys_menu"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统菜单_全量
// @Description	全量列表\n权限标识：无需权限（仅需登录）
// @Tags			SysMenuApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Success		200		{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_menu/all [get]
func All(c *fiber.Ctx) error {
	// 处理
	var args model.SysMenuDto
	var iSysMenuImpl = sysmenu.Impl()
	total, rows, err := iSysMenuImpl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp ListResp
	rp.Total = datacv.IntToStr(total)
	rp.Rows = make([]ListBody, len(rows))
	for idx, row := range rows {
		rp.Rows[idx] = ListBody{
			Id:         datacv.IntToStr(row.Id),
			Pid:        datacv.IntToStr(row.Pid),
			Types:      row.Types,
			Permis:     row.Permis,
			Remark:     row.Remark,
			CreateTime: datacv.IntToStr(row.CreateTime),
		}
	}
	return r.Resp(c, rp)
}
