package devprocess

import (
	"strings"

	"pkg/cronvalidate"
	"pkg/errs"
)

// DefaultCronExprEveryMinute 新建/开启定时未填写时默认（每分钟执行一次，与 robfig/cron 五段式一致）
const DefaultCronExprEveryMinute = "* * * * *"

func validateCronSettings(cronEnabled int8, cronExpr string) error {
	if cronEnabled != 0 && cronEnabled != 1 {
		return errs.New("定时启用状态无效")
	}
	if cronEnabled == 1 {
		s := strings.TrimSpace(cronExpr)
		if s == "" {
			return errs.New("定时执行须填写 Cron 表达式")
		}
		if err := cronvalidate.ValidateExpr(s); err != nil {
			return errs.New("Cron 表达式无效")
		}
	}
	return nil
}
