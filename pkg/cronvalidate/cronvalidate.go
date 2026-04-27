package cronvalidate

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// 仅接受标准 crontab 形式（可选秒字段），不接受 @daily、@every 等描述符。
var parser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

// ValidateExpr 校验 Cron 表达式（与 robfig/cron 解析规则一致）
func ValidateExpr(expr string) error {
	if expr == "" {
		return fmt.Errorf("表达式为空")
	}
	_, err := parser.Parse(expr)
	return err
}
