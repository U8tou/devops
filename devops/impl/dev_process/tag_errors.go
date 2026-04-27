package devprocess

import "pkg/errs"

var (
	errTagNameEmpty = errs.New("标签名称不能为空")
	errTagIds       = errs.New("存在无效的标签 id")
)
