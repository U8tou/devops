package mygit

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/ssh"
)

func TestRun(t *testing.T) {
	// 1. 拉取远程仓库至本地
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	localPath, err := GitFun(gitOpt)
	if err != nil {
		log.Fatalf("%v", err)
	}
	// 2. 执行构建命令
	cmds := strings.Split(runCmd, " ")
	c, err := LocalFun(cmds...)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("正在構建... \n", c)

	// 3. 调起ssh, 执行远程命令
	config := &ssh.ClientConfig{
		User: sshConf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshConf.Pwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", sshConf.Addr, config)
	if err != nil {
		log.Fatalf("%s %s \n", "Failed to dial:", err.Error())
	}
	defer client.Close()

	// 4. 上传文件
	// 5. 执行服务器命令
	// 6. 下载文件
	fmt.Println("正在執行遠程命令... \n", Ssh(SshOpt{}))

	// 刪除本地臨時文件
	err = tryRemoveAll(localPath, 3, 1*time.Second)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Println("臨時文件刪除成功:", localPath)
}
