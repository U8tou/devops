package mygit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// uploadFile 上传文件
func UploadFun(filePath string, client *ssh.Client) (shellMsg string, err error) {
	// 创建sftp客户端
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Failed to create SFTP client:", err.Error())
		return
	}
	defer sftpClient.Close()

	// 打开本地文件
	localFile, err := os.Open(filePath)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Unable to open local file:", err.Error())
		return
	}
	defer localFile.Close()
	fName := filepath.Base(localFile.Name())
	shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "File name to be uploaded:", fName)

	// 在SFTP服务器上创建并写入新文件
	remoteFile, err := sftpClient.Create("/root/" + fName)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Failed to create remote file:", err.Error())
		return
	}
	defer remoteFile.Close()

	// 将本地文件内容复制到远程文件
	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		shellMsg = fmt.Sprintf("%s %s %s \n", shellMsg, "Unable to write remote file:", err.Error())
		return
	}
	shellMsg = fmt.Sprintf("%s %s \n", shellMsg, "File uploaded successfully!!!")
	return
}
