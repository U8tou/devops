package devflow

import (
	"fmt"
	"strings"
)

// winCmdBuiltin 首词为 cmd 内建时不能无脑加 call（避免 if/set/for 等被改坏）。
// 以首词全小写精确匹配；ver 与 verify 不会混淆。
var winCmdBuiltin = map[string]struct{}{
	"set":      {},
	"setlocal": {}, "endlocal": {},
	"if": {}, "for": {}, "goto": {}, "else": {},
	"exit": {}, "chcp": {}, "title": {}, "color": {}, "pause": {},
	"echo": {}, "rem": {},
	"cd": {}, "chdir": {}, "pushd": {}, "popd": {},
	"path": {}, "call": {},
	"copy": {}, "del": {}, "dir": {},
	"md": {}, "mkdir": {}, "rmdir": {}, "rd": {},
	"ren": {}, "rename": {}, "move": {},
	"type": {}, "more": {},
	"ver": {}, "vol": {}, "verify": {},
	"shift": {}, "time": {}, "date": {},
	"cls": {}, "help": {}, "break": {},
}

func winBatchFirstToken(s string) string {
	s = strings.TrimLeft(s, " \t")
	if s == "" {
		return ""
	}
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			return s[:i]
		}
	}
	return s
}

func winBatchHasCompoundOp(line string) bool {
	for i := 0; i < len(line); i++ {
		if line[i] == '&' || line[i] == '|' {
			return true
		}
	}
	return false
}

// winBatchLineForFile 将一行用户脚本转为可写入 .cmd 的一行；对外部命令在默认情况下加 call。
// 第二返回值：false 仅当本行为空（略过）；否则 true（含空物理行在调用方另处理）。
func winBatchLineForFile(raw string) (string, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", false
	}
	low := strings.ToLower(s)
	if strings.HasPrefix(s, "#") {
		t := strings.TrimSpace(s[1:])
		if t == "" {
			return "rem", true
		}
		return "rem " + t, true
	}
	if winBatchHasCompoundOp(s) {
		return s, true
	}
	// 标签 :name（:: 为注释，见下）
	r := []rune(s)
	if len(r) >= 2 && r[0] == ':' && r[1] != ':' {
		return s, true
	}
	if len(r) > 0 && (r[0] == '@' || r[0] == '(' || r[0] == ')') {
		return s, true
	}
	// 注释/特殊
	if len(low) >= 2 && low[0] == ':' && low[1] == ':' {
		return s, true
	}
	// if( 不拆空格时
	if len(low) >= 3 && low[0] == 'i' && low[1] == 'f' && low[2] == '(' {
		return s, true
	}
	// path= 无空格
	if len(low) >= 5 && low[:4] == "path" && low[4] == '=' {
		return s, true
	}
	tok := strings.ToLower(winBatchFirstToken(s))
	if _, ok := winCmdBuiltin[tok]; ok {
		return s, true
	}
	// 批处理 echo. / echo, 为特殊 echo 形式
	if strings.HasPrefix(tok, "echo.") || strings.HasPrefix(tok, "echo,") {
		return s, true
	}
	// 外部 / 常由 .cmd 实现：call 回到父级脚本再继续（与手工 bat 推荐写法一致）
	return "call " + s, true
}

// winBatchFileFromScript 生成在 cmd 下单进程、同批处理上下文中执行的 .cmd 全文。
func winBatchFileFromScript(script string) (string, error) {
	if strings.TrimSpace(script) == "" {
		return "", fmt.Errorf("empty script")
	}
	normalized := strings.ReplaceAll(script, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")
	var b strings.Builder
	b.Grow(len(normalized) + 128)
	b.WriteString("@echo off\r\nsetlocal EnableExtensions\r\n")
	nSubstantive := 0
	for _, raw := range lines {
		line, ok := winBatchLineForFile(raw)
		if !ok {
			// 保留原脚本中的空行
			if strings.TrimSpace(raw) == "" {
				b.WriteString("\r\n")
			}
			continue
		}
		b.WriteString(line)
		b.WriteString("\r\n")
		nSubstantive++
	}
	if nSubstantive == 0 {
		return "", fmt.Errorf("empty script")
	}
	return b.String(), nil
}
