package sysrole

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_详情
// @Description	查询详情\n权限标识：sys:role:get
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		GetReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=GetResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_role/get [get]
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
	var iSysRoleImpl = sysrole.Impl()
	row, err := iSysRoleImpl.Get(c.Context(), datacv.StrToInt(rq.Id))
	if err != nil {
		return err
	}
	if row == nil {
		return r.Resp(c, nil)
	}
	// 响应
	rp := GetResp{
		Id:          datacv.IntToStr(row.Id),
		Name:        row.Name,
		Role:        row.Role,
		Status:      row.Status,
		MenuLinkage: row.MenuLinkage,
		DeptLinkage: row.DeptLinkage,
		Sort:        row.Sort,
		Remark:      row.Remark,
		MenuIds:     datacv.Int64SliceToStrSlice(row.Menus),
		DeptIds:     datacv.Int64SliceToStrSlice(row.Depts),
		CreateTime:  datacv.IntToStr(row.CreateTime),
		CreateBy:    datacv.IntToStr(row.CreateBy),
		UpdateTime:  datacv.IntToStr(row.UpdateTime),
		UpdateBy:    datacv.IntToStr(row.UpdateBy),
	}
	return r.Resp(c, rp)
}
