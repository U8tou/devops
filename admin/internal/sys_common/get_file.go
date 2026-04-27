package syscommon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pkg/conf"
	"pkg/errs"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

// allowedFileDirs 允许通过 get_file 访问的子目录（防止未授权枚举任意路径；如需开放更多目录可在此添加）
var allowedFileDirs = map[string]bool{
	"avatar": true, "export": true, "oper": true,
}

// @Summary		公共接口_获取文件
// @Description	获取文件（仅允许访问配置的 Local 目录下白名单子目录，防御路径穿越与目录枚举）
// @Tags			SysCommonApi
// @Produce		octet-stream
// @Param			dir	path		string		true	"文件夹"
// @Param			obj	path		string		true	"文件名"
// @Success		200	{file}		file		"文件内容"
// @Router			/common/get_file/{dir}/{obj} [get]
func GetFile(c *fiber.Ctx) error {
	dir := strings.TrimSpace(c.Params("dir"))
	obj := strings.TrimSpace(c.Params("obj"))
	if dir == "" || obj == "" {
		return errs.Args(errors.New("文件夹或文件名不能为空"))
	}
	if !allowedFileDirs[dir] {
		return errs.Args(errors.New("不允许访问该目录"))
	}
	basePath, err := filepath.Abs(conf.File.Local)
	if err != nil {
		return errs.Sys(err)
	}
	f := filepath.Join(basePath, filepath.Clean(filepath.Join(dir, obj)))
	absPath, err := filepath.Abs(f)
	if err != nil {
		return errs.Sys(err)
	}
	rel, err := filepath.Rel(basePath, absPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return errs.Args(errors.New("非法路径"))
	}
	file, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	// 防止 Content-Disposition 注入：仅保留安全字符作为 filename
	safeName := sanitizeContentDispositionFilename(filepath.Base(obj))
	if safeName == "" {
		safeName = "download"
	}
	c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, safeName))
	return c.Send(file)
}

// sanitizeContentDispositionFilename 去除换行、双引号等，防止响应头注入
func sanitizeContentDispositionFilename(name string) string {
	var b strings.Builder
	for _, r := range name {
		if r == '"' || r == '\r' || r == '\n' || r == '\\' || !unicode.IsPrint(r) {
			continue
		}
		b.WriteRune(r)
	}
	s := b.String()
	if len(s) > 255 {
		s = s[:255]
	}
	return s
}
