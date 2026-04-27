package sysmenu

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysmenu "system/impl/sys_menu"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统菜单_详情
// @Description	查询详情\n权限标识：sys:menu:get
// @Tags			SysMenuApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		GetReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=GetResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_menu/get [get]
func Get(c *fiber.Ctx) error {
	// 请求
	var rq GetReq
	if err := c.QueryParser(&rq); err != nil {
		return err
	}
	err := validator.Struct(rq)
	if err != nil {
		return err
	}
	// 处理
	var iSysMenuImpl = sysmenu.Impl()
	row, err := iSysMenuImpl.Get(c.Context(), datacv.StrToInt(rq.Id))
	if err != nil {
		return err
	}
	if row == nil {
		return r.Resp(c, nil)
	}
	// 响应
	rp := GetResp{
		Id:         datacv.IntToStr(row.Id),
		Pid:        datacv.IntToStr(row.Pid),
		Types:      row.Types,
		Permis:     row.Permis,
		Remark:     row.Remark,
		CreateTime: datacv.IntToStr(row.CreateTime),
		CreateBy:   datacv.IntToStr(row.CreateBy),
		UpdateTime: datacv.IntToStr(row.UpdateTime),
		UpdateBy:   datacv.IntToStr(row.UpdateBy),
	}
	return r.Resp(c, rp)
}
