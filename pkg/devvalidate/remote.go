package devvalidate

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// ParseRemoteHostPort 解析 user@host:port、host:port、host，默认 SSH 端口 22
func ParseRemoteHostPort(s string) (host, port string, err error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", "", errors.New("远程地址不能为空")
	}
	hostPart := s
	if i := strings.LastIndex(s, "@"); i >= 0 {
		hostPart = strings.TrimSpace(s[i+1:])
	}
	if hostPart == "" {
		return "", "", errors.New("远程地址格式无效")
	}
	if h, p, e := net.SplitHostPort(hostPart); e == nil {
		return h, p, nil
	}
	return hostPart, "22", nil
}

// ValidateRemoteTCP 通过 TCP 拨测校验远程主机端口是否可达
func ValidateRemoteTCP(ctx context.Context, address string) error {
	address = strings.TrimSpace(address)
	host, port, err := ParseRemoteHostPort(address)
	if err != nil {
		return err
	}
	d := net.Dialer{Timeout: 10 * time.Second}
	addr := net.JoinHostPort(host, port)
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("远程不可达: %v", err)
	}
	_ = conn.Close()
	return nil
}
