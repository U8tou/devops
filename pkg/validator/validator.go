package validator

import (
	"github.com/go-playground/validator/v10"
)

// DefaultValidator 全局单例，复用以减少分配、利用结构体缓存
var DefaultValidator = validator.New()

// Struct 校验结构体，等价于 DefaultValidator.Struct(s)
func Struct(s interface{}) error {
	return DefaultValidator.Struct(s)
}
