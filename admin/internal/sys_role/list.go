package sysrole

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"
	"system/model"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_分页
// @Description	分页列表\n权限标识：sys:role:list
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		json
// @Param			args	query		ListReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=ListResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_role/list [get]
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
	var args model.SysRoleDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.Name = rq.Name
	args.Role = rq.Role
	args.Status = rq.Status
	var iSysRoleImpl = sysrole.Impl()
	total, rows, err := iSysRoleImpl.List(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp ListResp
	rp.Total = datacv.IntToStr(total)
	rp.Rows = make([]ListBody, len(rows))
	for idx, row := range rows {
		rp.Rows[idx] = ListBody{
			Id:          datacv.IntToStr(row.Id),
			Name:        row.Name,
			Role:        row.Role,
			Status:      row.Status,
			MenuLinkage: row.MenuLinkage,
			DeptLinkage: row.DeptLinkage,
			Sort:        row.Sort,
			Remark:      row.Remark,
			CreateTime:  datacv.IntToStr(row.CreateTime),
			UpdateTime:  datacv.IntToStr(row.UpdateTime),
		}
	}
	return r.Resp(c, rp)
}
