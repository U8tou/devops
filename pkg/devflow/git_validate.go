package devflow

import (
	"context"
	"time"
)

const gitValidateTimeout = 30 * time.Second

// ValidateGitRepo 使用 go-git 列远程引用校验仓库与分支（支持 HTTPS Basic / SSH 私钥，见 gitAuthType）
func ValidateGitRepo(ctx context.Context, params map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, gitValidateTimeout)
	defer cancel()
	return ValidateGitRepoWithGoGit(ctx, params)
}
