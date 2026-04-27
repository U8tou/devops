package sysuser

import (
	"bytes"
	"strconv"

	"admin/internal/datapermctx"
	"pkg/tools/datacv"
	"pkg/tools/excel"
	sysuser "system/impl/sys_user"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// @Summary		用户导出
// @Description	按条件导出用户列表为 Excel（仅基础信息，不含密码）\n权限标识：sys:user:list
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Accept			x-www-form-urlencoded
// @Produce		octet-stream
// @Param			args	query		ListReq	true	"与 list 接口一致的筛选参数"
// @Success		200		{file}		file	"用户列表.xlsx"
// @Router			/sys_user/export [get]
func Export(c *fiber.Ctx) error {
	var rq ListReq
	if err := c.QueryParser(&rq); err != nil {
		return err
	}
	if rq.Size <= 0 || rq.Size > 100000 {
		rq.Size = 10000
	}
	if rq.Current <= 0 {
		rq.Current = 1
	}
	var args model.SysUserDto
	args.Current = rq.Current
	args.Size = rq.Size
	args.Id = datacv.StrToInt(rq.Id)
	args.UserName = rq.UserName
	args.NickName = rq.NickName
	args.Email = rq.Email
	args.Phone = rq.Phone
	args.Sex = rq.Sex
	args.Status = rq.Status
	args.DeptId = datacv.StrToInt(rq.DeptId)
	if err := datapermctx.ApplyUserDto(c, &args); err != nil {
		return err
	}

	_, rows, err := sysuser.Impl().List(c.Context(), &args)
	if err != nil {
		return err
	}

	// 导出不含密码列，与 UserExportHeaders 顺序一致
	excelRows := make([][]string, 0, len(rows))
	for _, row := range rows {
		excelRows = append(excelRows, []string{
			row.UserName,
			row.NickName,
			row.UserType,
			row.Email,
			row.PhoneArea,
			row.Phone,
			strconv.FormatInt(int64(row.Sex), 10),
			strconv.FormatInt(int64(row.Status), 10),
			row.Address,
			row.Remark,
		})
	}

	f, err := excel.WriteSheet(UserExportHeaders, excelRows)
	if err != nil {
		return err
	}
	if err := excel.StyleHeaderRow(f, "Sheet1", len(UserExportHeaders), userExportColWidths); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", `attachment; filename="用户列表.xlsx"`)
	return c.Send(buf.Bytes())
}
