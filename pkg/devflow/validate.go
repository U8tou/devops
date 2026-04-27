package devflow

import (
	"context"
	"fmt"
	"strings"
)

// ValidateNode 校验节点（kind 与前端 data.kind 一致）
func ValidateNode(ctx context.Context, kind string, params map[string]interface{}) error {
	k := strings.TrimSpace(strings.ToLower(kind))
	switch k {
	case "git_repo":
		return ValidateGitRepo(ctx, params)
	case "ssh_connection":
		return ValidateSSHConnectionParams(ctx, params)
	case "upload_servers":
		// 与 ssh_connection 相同结构化字段时做完整 SSH 校验
		if strings.TrimSpace(strParam(params, "username")) != "" {
			return ValidateSSHConnectionParams(ctx, params)
		}
		host := strings.TrimSpace(strParam(params, "host"))
		if host == "" {
			return fmt.Errorf("host required")
		}
		return ValidateUploadServersHost(ctx, host)
	default:
		return fmt.Errorf("unsupported kind: %s", kind)
	}
}
