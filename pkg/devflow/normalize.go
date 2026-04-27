package devflow

import "strings"

// NormalizeNodeKind 与前端 normalizeNodeKind 对齐
func NormalizeNodeKind(kind string) string {
	k := strings.TrimSpace(kind)
	if k == "build" || k == "remote_script" || k == "start_service" {
		return "execute_script"
	}
	switch k {
	case "git_repo", "ssh_connection", "remote_ssh_script", "execute_script", "upload_servers", "remote_download":
		return k
	default:
		return "git_repo"
	}
}
