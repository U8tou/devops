package mygit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// downloadFile 下载文件
func DownloadFun(filePath string, client *ssh.Client) (shellMsg string, err error) {
	// 创建sftp客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Failed to create SFTP client:", err.Error())
		return
	}
	defer sftpClient.Close()
	// 在SFTP服务器上打开文件
	remoteFile, err := sftpClient.Open(filePath)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Failed to open remote file:", err.Error())
		return
	}
	defer remoteFile.Close()

	// 创建本地文件
	localFile, err := os.Create(filepath.Base(remoteFile.Name()))
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Failed to create local file:", err.Error())
		return
	}
	defer localFile.Close()
	// 将远程文件内容复制到本地文件
	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Unable to write local file:", err.Error())
		return
	}
	shellMsg = fmt.Sprintf("%s %s \n", shellMsg, "File download successfully!!!")
	return
}
