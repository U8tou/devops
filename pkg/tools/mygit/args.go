package mygit

type GitOpt struct {
	Dir           string // 臨時文件夾
	RepoURL       string // 遠程git倉庫
	ReferenceName string // 獲取分支
	User          string // 用戶名
	Pwd           string // 密碼
}

type SshOpt struct {
	Addr string
	User string
	Pwd  string
}

type ProcessOpt struct {
	SshList  []SshOpt  // 远程地址, 若为空，则无远程操作（3_上传文件 4_执行ssh命令 5_下载文件 不可用）
	OperList []OperOpt // 执行操作
}

// 操作
type OperOpt struct {
	Typ        int    // 1_拉取代码 2_执行本地命令 3_上传文件 4_执行ssh命令 5_下载文件
	Cmd        string // 命令
	LocalDir   string // 本地文件夹
	LocalFile  string // 本地文件
	RemoteDir  string // 远程文件夹
	RemoteFile string // 远程文件
	GitConf    GitOpt // git仓库
}

var (
	gitOpt = GitOpt{
		RepoURL:       "https://gitee.com/U8tou/bigo_server.git",
		ReferenceName: "master",
		Dir:           "E:\\tmp",
		User:          "u8tou",
		Pwd:           "QWERasd0",
	}
	sshConf = SshOpt{
		Addr: "47.119.132.182:22",
		User: "root",
		Pwd:  "QWERasd0..",
	}
	runCmd = `
java -version && pnpm -v && npm -v
`
	processList = ProcessOpt{
		SshList:  []SshOpt{},
		OperList: []OperOpt{},
	}
)
