package sysrole

import (
	r "pkg/resp"
	"pkg/tools/datacv"
	sysrole "system/impl/sys_role"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统角色_全量
// @Description	获取所有角色\n权限标识：无需权限（仅需登录）
// @Tags			SysRoleApi
// @Security		ApiKeyAuth
// @Produce		json
// @Success		200	{object}	r.MyResp{data=[]ListBody}	"Success"
// @Failure		500	{object}	r.MyResp						"Failure"
// @Router			/sys_role/all [get]
func All(c *fiber.Ctx) error {
	// 处理
	var iSysRoleImpl = sysrole.Impl()
	rows, err := iSysRoleImpl.All(c.Context())
	if err != nil {
		return err
	}
	// 响应
	result := make([]ListBody, len(rows))
	for idx, row := range rows {
		result[idx] = ListBody{
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
	return r.Resp(c, result)
}
