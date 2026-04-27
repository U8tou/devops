package conf

import "strings"

// FileUrl 将文件路径转为可访问的完整 URL。
// 若为网络图片（http/https）则直接返回；非网络图片则返回 BaseUrl + /uploads/ + path。
func FileUrl(path string) string {
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return path
	// base := strings.TrimSuffix(File.BaseUrl, "/")
	// return base + "/uploads/" + strings.TrimPrefix(path, "/")
}
