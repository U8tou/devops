package sysmenu

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysmenu "system/impl/sys_menu"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统菜单_分页
// @Description	分页列表\n权限标识：sys:menu:list
// @Tags			SysMenuApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		ListReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_menu/list [get]
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
	var args model.SysMenuDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.Pid = datacv.StrToInt(rq.Pid)
	args.Types = rq.Types
	args.Permis = rq.Permis
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
