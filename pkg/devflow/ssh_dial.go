package devflow

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const sshDialTimeout = 15 * time.Second

// DialSSH 建立 SSH 客户端（密码或私钥）；addr 为 host:port
func DialSSH(ctx context.Context, addr, username, password, privateKeyPEM, authType string) (*ssh.Client, error) {
	authType = strings.TrimSpace(strings.ToLower(authType))
	if authType == "" {
		authType = "key"
	}
	var auth []ssh.AuthMethod
	switch authType {
	case "password":
		if password == "" {
			return nil, fmt.Errorf("password auth requires password")
		}
		auth = []ssh.AuthMethod{ssh.Password(password)}
	default:
		if strings.TrimSpace(privateKeyPEM) == "" {
			return nil, fmt.Errorf("key auth requires private key content")
		}
		signer, err := ssh.ParsePrivateKey([]byte(privateKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}

	cfg := &ssh.ClientConfig{
		User:            username,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         sshDialTimeout,
	}

	d := net.Dialer{Timeout: sshDialTimeout}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcp dial: %w", err)
	}
	c, chans, reqs, err := ssh.NewClientConn(conn, addr, cfg)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ssh handshake: %w", err)
	}
	return ssh.NewClient(c, chans, reqs), nil
}

// ValidateSSHConnectionParams 校验 ssh_connection 节点参数（map 来自 JSON）
func ValidateSSHConnectionParams(ctx context.Context, params map[string]interface{}) error {
	host := strParam(params, "host")
	port := strParam(params, "port")
	if strings.TrimSpace(port) == "" {
		port = "22"
	}
	user := strParam(params, "username")
	authType := strParam(params, "authType")
	if strings.TrimSpace(host) == "" {
		return fmt.Errorf("host required")
	}
	if strings.TrimSpace(user) == "" {
		return fmt.Errorf("username required")
	}
	addr := net.JoinHostPort(strings.TrimSpace(host), strings.TrimSpace(port))
	pwd := strParam(params, "password")
	key := strParam(params, "privateKey")
	cli, err := DialSSH(ctx, addr, user, pwd, key, authType)
	if err != nil {
		return err
	}
	_ = cli.Close()
	return nil
}

// ValidateUploadServersHost 仅 SSH 连通性（与前端「校验」一致）
func ValidateUploadServersHost(ctx context.Context, hostLine string) error {
	t, ok := ParseUploadHost(hostLine)
	if !ok || strings.TrimSpace(t.Host) == "" {
		return fmt.Errorf("invalid upload host")
	}
	port := t.Port
	if port == "" {
		port = "22"
	}
	user := t.Username
	if user == "" {
		user = "root"
	}
	// 上传节点校验无密码：仅 TCP + SSH 协议握手无法在无凭据下完成；此处做 TCP 连通性
	addr := net.JoinHostPort(t.Host, port)
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("tcp dial %s: %w", addr, err)
	}
	_ = conn.Close()
	return nil
}
