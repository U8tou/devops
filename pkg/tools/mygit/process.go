package mygit

import "sync"

// Process 执行流程
func Process(pro ProcessOpt) {
	var wg sync.WaitGroup
	for _, ssh := range pro.SshList {
		wg.Add(1)
		println(ssh.Addr)
		go sshFun(&wg)
	}
	// 等待所有goroutine完成
	wg.Wait()
	// 1_拉取代码 2_执行本地命令 3_上传文件 4_执行ssh命令 5_下载文件
	for _, oper := range pro.OperList {
		switch oper.Typ {
		case 1:
			// 拉取Git仓库代码
		case 2:
			// 执行本地命令
		case 3:
			// 上传文件
		case 4:
			// 执行ssh命令
		case 5:
			// 下载文件
		}

	}
}

func sshFun(wg *sync.WaitGroup) {
	defer wg.Done()
}
