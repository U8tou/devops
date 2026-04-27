package devflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport"
	githttp "github.com/go-git/go-git/v6/plumbing/transport/http"
	gitssh "github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"github.com/go-git/go-git/v6/storage/memory"
	"golang.org/x/crypto/ssh"
)

func strParamLower(m map[string]interface{}, key string) string {
	return strings.ToLower(strings.TrimSpace(strParam(m, key)))
}

// buildGitAuth 根据 params（gitAuthType + 凭证字段）与仓库 URL 构造 go-git Auth；不记录敏感信息。
func buildGitAuth(repoURL string, params map[string]interface{}) (transport.AuthMethod, error) {
	authKind := strParamLower(params, "gitAuthType")
	if authKind == "" || authKind == "none" {
		return nil, nil
	}
	isHTTP := strings.HasPrefix(repoURL, "http://") || strings.HasPrefix(repoURL, "https://")
	isSSH := strings.HasPrefix(repoURL, "git@") || strings.HasPrefix(repoURL, "ssh://")

	switch authKind {
	case "http":
		user := strParam(params, "httpUsername")
		pass := strParam(params, "httpPassword")
		if !isHTTP {
			return nil, fmt.Errorf("HTTPS 凭证仅适用于 http(s) 仓库地址")
		}
		if user == "" && pass == "" {
			return nil, fmt.Errorf("请填写 HTTPS 用户名或密码/令牌")
		}
		return &githttp.BasicAuth{Username: user, Password: pass}, nil
	case "ssh_key", "sshkey":
		keyPEM := strParam(params, "sshPrivateKey")
		if keyPEM == "" {
			return nil, fmt.Errorf("请填写 SSH 私钥（PEM）")
		}
		if !isSSH {
			return nil, fmt.Errorf("SSH 私钥仅适用于 git@ 或 ssh:// 仓库地址")
		}
		pub, err := gitssh.NewPublicKeys("git", []byte(keyPEM), "")
		if err != nil {
			return nil, fmt.Errorf("parse SSH private key: %w", err)
		}
		pub.HostKeyCallback = ssh.InsecureIgnoreHostKey()
		return pub, nil
	default:
		return nil, nil
	}
}

// ValidateGitRepoWithGoGit 使用 go-git 列远程引用校验分支可达（支持 HTTP/SSH 凭证）
func ValidateGitRepoWithGoGit(ctx context.Context, params map[string]interface{}) error {
	repo := strings.TrimSpace(strParam(params, "repositoryUrl"))
	branch := strings.TrimSpace(strParam(params, "branch"))
	if repo == "" {
		return fmt.Errorf("repositoryUrl required")
	}
	if branch == "" {
		branch = "main"
	}
	auth, err := buildGitAuth(repo, params)
	if err != nil {
		return err
	}
	stor := memory.NewStorage()
	rem := git.NewRemote(stor, &config.RemoteConfig{Name: "origin", URLs: []string{repo}})
	refs, err := rem.ListContext(ctx, &git.ListOptions{Auth: auth})
	if err != nil {
		return fmt.Errorf("list remote refs: %w", err)
	}
	want := plumbing.NewBranchReferenceName(branch).String()
	for _, rf := range refs {
		if rf.Name().String() == want {
			return nil
		}
	}
	return fmt.Errorf("branch %q not found on remote", branch)
}

// CloneGitRepoWithGoGit 使用 go-git 克隆到 dest
func CloneGitRepoWithGoGit(ctx context.Context, params map[string]interface{}, dest string) error {
	repo := strings.TrimSpace(strParam(params, "repositoryUrl"))
	branch := strings.TrimSpace(strParam(params, "branch"))
	if branch == "" {
		branch = "main"
	}
	if repo == "" {
		return fmt.Errorf("repositoryUrl required")
	}
	auth, err := buildGitAuth(repo, params)
	if err != nil {
		return err
	}
	opts := &git.CloneOptions{
		URL:             repo,
		Auth:            auth,
		SingleBranch:    true,
		ReferenceName:   plumbing.NewBranchReferenceName(branch),
		NoCheckout:      false,
		InsecureSkipTLS: false,
	}
	if boolParam(params, "shallowClone") {
		opts.Depth = 1
	}
	_, err = git.PlainCloneContext(ctx, dest, opts)
	if err != nil {
		return fmt.Errorf("git clone: %w", err)
	}
	return nil
}
