package mygit

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 執行構建命令
func LocalFun(arg ...string) (string, error) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// Windows 需要通过 cmd.exe 执行内置命令，/c 表示执行完后关闭
		args := append([]string{"/C"}, arg...)
		cmd = exec.Command("cmd", args...)
	case "linux", "darwin":
		// Unix-like 系统直接调用命令名
		args := append([]string{}, arg...)
		cmd = exec.Command("bash", args...)
	default:
		return "", fmt.Errorf("不支持的系统: %s\n", runtime.GOOS)
	}

	// 继承父进程所有环境变量(否則可能找不到環境)
	cmd.Env = os.Environ()
	// 设置环境变量
	// cmd.Env = append(cmd.Env, "MY_GO_ENV=production")

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		out, _ := fmtPrint(output)
		return out, fmt.Errorf("执行命令失败: %v\n", err)
	}
	return fmtPrint(output)
}

func fmtPrint(output []byte) (string, error) {
	// 如果是 Windows 系统，将 GBK 输出转换为 UTF-8
	if runtime.GOOS == "windows" {
		// 使用 simplifiedchinese.GBK.NewDecoder() 创建一个 GBK 解码器
		decoder := simplifiedchinese.GBK.NewDecoder()
		// 将 GBK 编码的字节切片转换为 UTF-8 编码
		utf8Output, _, err := transform.Bytes(decoder, output)
		if err != nil {
			return string(output), fmt.Errorf("转码失败: %v\n", err)
		}
		return string(utf8Output), err
	} else {
		// 在非 Windows 系统上，输出通常已经是 UTF-8，可以直接打印
		return string(output), nil
	}
}
