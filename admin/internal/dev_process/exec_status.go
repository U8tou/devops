package dev_process

import "strings"

// normalizeListExecStatus 将库中的 last_exec_result 转为列表展示用状态码。
// 新版存 success/failed/cancelled；旧版曾把整段日志存在该列，这里做兼容推断。
func normalizeListExecStatus(raw string) string {
	s := strings.TrimSpace(raw)
	switch s {
	case "success", "failed", "cancelled":
		return s
	}
	if s == "" {
		return ""
	}
	if strings.Contains(s, "\nerror:") {
		return "failed"
	}
	return "success"
}
