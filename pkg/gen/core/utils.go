package core

import (
	"reflect"
	"strings"
	"unicode"
)

// toLowerCamelCase 将大驼峰命名的字符串转换为小驼峰
func toLowerCamelCase(s string) string {
	firstChar := strings.ToLower(s[:1])
	rest := s[1:]
	return firstChar + rest
}

// pascalToKebab 将PascalCase转换为kebab-case，如 Post→post, OrderItem→order-item
func pascalToKebab(s string) string {
	return strings.ReplaceAll(camelCaseToSnakeCase(s), "_", "-")
}

// camelCaseToSnakeCase 将大驼峰命名的字符串转换为下划线命名
func camelCaseToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	// 定义一个缓冲区来存储结果
	var buffer strings.Builder
	// 遍历字符串的每个字符
	for i, c := range s {
		// 如果是大写字母，并且不是字符串的第一个字符
		if unicode.IsUpper(c) && i != 0 {
			// 在缓冲区中写入下划线
			err := buffer.WriteByte('_')
			if err != nil {
				panic(err)
			}
		}
		// 将字符转换为小写并写入缓冲区
		_, err := buffer.WriteRune(unicode.ToLower(c))
		if err != nil {
			panic(err)
		}
	}
	// 返回缓冲区的字符串表示
	return buffer.String()
}

func GetStructName(v any) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// 获取两个符号之间的内容
func GetBetweenSub(str string, startStr string, endStar string) string {
	start := strings.Index(str, startStr)
	if start == -1 {
		return ""
	}
	end := strings.Index(str[start+1:], endStar)
	if end == -1 {
		return ""
	}
	return str[start+1 : start+1+end]
}
