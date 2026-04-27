package mygit

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/jaevor/go-nanoid"
)

func GitFun(bui GitOpt) (string, error) {
	gen, err := nanoid.Canonic()
	if err != nil {
		return "", err
	}
	// 替换成本地保存路径
	localPath := filepath.Join(bui.Dir, gen())
	log.Println("localPath:", localPath)

	auth := &http.BasicAuth{
		Username: bui.User,
		Password: bui.Pwd,
	}

	// 1. 克隆仓库 (第一次拉取)
	// 替换成你的仓库 URL
	repoURL := bui.RepoURL
	_, err = git.PlainClone(localPath, &git.CloneOptions{
		URL:           repoURL,
		Auth:          auth,                                               // 如果需要认证
		ReferenceName: plumbing.NewBranchReferenceName(bui.ReferenceName), // 指定分支
	})
	if err != nil {
		return "", err
	}
	fmt.Println("克隆成功:", localPath)

	return localPath, nil
}

func tryRemoveAll(path string, retries int, delay time.Duration) error {
	for i := range retries {
		err := os.RemoveAll(path)
		if err == nil {
			return nil // 成功删除
		}
		fmt.Printf("删除失败，第 %d 次重试，等待 %v: %v\n", i+1, delay, err)
		time.Sleep(delay)
	}
	return fmt.Errorf("无法在 %d 次尝试后删除目录 %s", retries, path)
}
