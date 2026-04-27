package devflow

import (
	"strings"
)

// SshTarget 与前端 upload_servers host 解析一致（user@host:port）
type SshTarget struct {
	Username string
	Host     string
	Port     string
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// ParseUploadHost 解析「上传部署」的 host 字段
func ParseUploadHost(hostRaw string) (SshTarget, bool) {
	s := strings.TrimSpace(hostRaw)
	if s == "" {
		return SshTarget{}, false
	}
	at := strings.IndexByte(s, '@')
	if at < 0 {
		lastColon := strings.LastIndex(s, ":")
		if lastColon > 0 && isDigits(s[lastColon+1:]) {
			return SshTarget{Username: "", Host: s[:lastColon], Port: s[lastColon+1:]}, true
		}
		return SshTarget{Username: "", Host: s, Port: "22"}, true
	}
	user := s[:at]
	rest := s[at+1:]
	lastColon := strings.LastIndex(rest, ":")
	if lastColon > 0 && len(rest) > lastColon+1 && isDigits(rest[lastColon+1:]) {
		return SshTarget{Username: user, Host: rest[:lastColon], Port: rest[lastColon+1:]}, true
	}
	return SshTarget{Username: user, Host: rest, Port: "22"}, true
}
