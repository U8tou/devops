package sysuser

import (
	"bytes"
	"pkg/tools/excel"

	"github.com/gofiber/fiber/v2"
)

// 用户导入/导出 Excel 表头（仅基础信息），与模板、导入共用（含密码列）
var UserExcelHeaders = []string{
	"用户账号", "用户昵称", "用户类型", "邮箱", "电话区号", "电话号码",
	"性别", "状态", "联系地址", "备注", "密码",
}

// 导出用表头：不含密码列
var UserExportHeaders = []string{
	"用户账号", "用户昵称", "用户类型", "邮箱", "电话区号", "电话号码",
	"性别", "状态", "联系地址", "备注",
}

// 模板示例行（可选，便于用户理解格式）
var userTemplateExampleRow = []string{
	"zhangsan", "张三", "00", "zhangsan@example.com", "+86", "13800138000",
	"1", "1", "北京市朝阳区", "示例备注", "123456",
}

// 模板列宽（与 UserExcelHeaders 一一对应，保证表头与示例内容完整显示）
var userTemplateColWidths = []float64{
	14, 14, 10, 24, 12, 14, 8, 8, 28, 18, 12, // 用户账号~密码
}

// 导出列宽（与 UserExportHeaders 一一对应，不含密码列）
var userExportColWidths = []float64{
	14, 14, 10, 24, 12, 14, 8, 8, 28, 18, // 用户账号~备注
}

// @Summary		用户导入模板下载
// @Description	下载用户导入 Excel 模板（仅含基础信息列）\n权限标识：sys:user:list
// @Tags			SysUserApi
// @Security		ApiKeyAuth
// @Produce		octet-stream
// @Success		200	{file}	file	"用户导入模板.xlsx"
// @Router			/sys_user/template [get]
func Template(c *fiber.Ctx) error {
	rows := [][]string{userTemplateExampleRow}
	f, err := excel.WriteSheet(UserExcelHeaders, rows)
	if err != nil {
		return err
	}
	if err := excel.StyleHeaderRow(f, "Sheet1", len(UserExcelHeaders), userTemplateColWidths); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return err
	}
	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", `attachment; filename="用户导入模板.xlsx"`)
	return c.Send(buf.Bytes())
}
