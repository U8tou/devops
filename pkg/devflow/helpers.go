package devflow

import (
	"fmt"
)

// formatByteSize 人类可读字节大小（二进制 KiB/MiB/GiB）
func formatByteSize(n int64) string {
	if n < 0 {
		n = 0
	}
	if n < 1024 {
		return fmt.Sprintf("%d B", n)
	}
	// 先换算到 KiB，再向 MiB 推进；若从 n 开始除且 u 从 0 递增，会错用 MiB 表示 1KiB
	xf := float64(n) / 1024
	u := 0
	units := []string{"KiB", "MiB", "GiB", "TiB"}
	for xf >= 1024 && u < len(units)-1 {
		xf /= 1024
		u++
	}
	return fmt.Sprintf("%.2f %s", xf, units[u])
}

func strParam(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%.0f", t)
	case bool:
		if t {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprint(t)
	}
}

func boolParam(m map[string]interface{}, key string) bool {
	if m == nil {
		return false
	}
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t == "true" || t == "1"
	case float64:
		return t != 0
	default:
		return false
	}
}
