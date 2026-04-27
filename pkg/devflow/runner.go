package devflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"
)

// cappedMergeWriter 供 SSH / exec 的 Stdout+Stderr 安全合并，且单步总字节有上限，避免大输出撑爆内存或日志。
type cappedMergeWriter struct {
	mu        sync.Mutex
	limit     int
	buf       strings.Builder
	truncated bool
}

func (w *cappedMergeWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.limit <= 0 {
		w.truncated = len(p) > 0
		return len(p), nil
	}
	remaining := w.limit - w.buf.Len()
	if remaining <= 0 {
		if len(p) > 0 {
			w.truncated = true
		}
		return len(p), nil
	}
	if len(p) > remaining {
		_, _ = w.buf.Write(p[:remaining])
		w.truncated = true
		return len(p), nil
	}
	return w.buf.Write(p)
}

func (w *cappedMergeWriter) resultString() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	s := w.buf.String()
	if w.truncated {
		return s + fmt.Sprintf("\n...(subprocess output truncated, max %d bytes per step)", w.limit)
	}
	return s
}

// 各常量：单步子进程/单步远程 SSH 合并输出、聚合日志、落库 分层限制，防脚本刷爆库与页面。
const (
	maxScriptIOBytes  = 48_000  // 单步本地/SSH 子进程收集的原始输出上限
	maxRunLogBytes    = 100_000 // 整个 Run 的聚合文本上限（行拼接后）
	maxLogLineBytes   = 7_200   // 单条 append 到聚合日志的字符串上限（防单行超大）
)

// RunProgress 可选：供 SSE 等场景推送节点与日志行；字段均可为空。
// transferBytes：上传/下载类节点传输的字节总数，无则 0。
// skipped 为 true 表示本节点因「流程开关」关闭而未执行，duration/transfer 均为 0。
type RunProgress struct {
	OnLog       func(line string)
	OnNodeStart func(nodeID, kind string, startedAt time.Time)
	OnNodeEnd   func(nodeID, kind string, err error, duration time.Duration, transferBytes int64, skipped bool)
}

// Run 按拓扑序同步执行流程节点，返回聚合日志（由 BuildLastExecRecord 写入 last_exec_log / last_exec_result）。
// procEnv 为流程级环境变量，会并入本地 exec/git 子进程环境（同名覆盖）。
// progress 非 nil 时同步回调日志行与节点起止（err 非 nil 表示该节点失败）。
func Run(ctx context.Context, flowJSON, workspaceRoot, processID string, procEnv map[string]string, progress *RunProgress) (string, error) {
	if strings.TrimSpace(workspaceRoot) == "" {
		return "", fmt.Errorf("workspace root is empty")
	}
	nodes, edges, err := ParseFlow(flowJSON)
	if err != nil {
		return "", fmt.Errorf("parse flow: %w", err)
	}
	ordered := TopoOrderedNodes(nodes, edges)
	skip := BuildDownstreamFlowSkip(nodes, edges)
	baseDir := filepath.Join(workspaceRoot, "dev_process", processID)
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return "", fmt.Errorf("mkdir workspace: %w", err)
	}

	var b strings.Builder
	var runLogFull bool
	const furtherOmitted = "...(further log omitted)"
	appendLog := func(line string) {
		if runLogFull {
			return
		}
		if len(line) > maxLogLineBytes {
			line = line[:maxLogLineBytes-24] + "\n...(line truncated)"
		}
		add := len(line)
		if b.Len() > 0 {
			add++
		}
		if b.Len()+add > maxRunLogBytes {
			room := maxRunLogBytes - b.Len()
			if b.Len() > 0 {
				room-- // 换行
			}
			if room < 24 {
				if b.Len() > 0 {
					b.WriteByte('\n')
				}
				b.WriteString(furtherOmitted)
				runLogFull = true
				if progress != nil && progress.OnLog != nil {
					progress.OnLog(furtherOmitted)
				}
				return
			}
			piece := line
			if len(piece) > room-20 {
				piece = piece[:room-20] + "\n...(truncated)"
			}
			if b.Len() > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(piece)
			runLogFull = true
			if progress != nil && progress.OnLog != nil {
				progress.OnLog(piece)
			}
			return
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(line)
		if progress != nil && progress.OnLog != nil {
			progress.OnLog(line)
		}
	}

	for _, n := range ordered {
		kind, params, err := ParseNodeData(n.Data)
		if err != nil {
			return b.String(), fmt.Errorf("node %s: parse data: %w", n.ID, err)
		}
		logName := NodeDisplayNameForLog(n.Data, kind)
		linePrefix := fmt.Sprintf("[%s] %s", n.ID, logName)
		if _, off := skip[n.ID]; off {
			stepStarted := time.Now()
			if progress != nil && progress.OnNodeStart != nil {
				progress.OnNodeStart(n.ID, kind, stepStarted)
			}
			appendLog(fmt.Sprintf("%s: skipped (flow off)", linePrefix))
			if progress != nil && progress.OnNodeEnd != nil {
				progress.OnNodeEnd(n.ID, kind, nil, 0, 0, true)
			}
			continue
		}
		stepStarted := time.Now()
		if progress != nil && progress.OnNodeStart != nil {
			progress.OnNodeStart(n.ID, kind, stepStarted)
		}
		appendLog(fmt.Sprintf("%s: start (%s)", linePrefix, stepStarted.Format("2006-01-02 15:04:05.000")))
		transferBytes, stepErr := runNodeStep(ctx, kind, params, n.ID, nodes, edges, baseDir, procEnv, appendLog)
		elapsed := time.Since(stepStarted)
		transferSuffix := ""
		if transferBytes > 0 {
			transferSuffix = "，传输 " + formatByteSize(transferBytes)
		}
		if stepErr != nil {
			appendLog(fmt.Sprintf("%s: %s (耗时 %s%s)", linePrefix, stepErr.Error(), formatNodeStepDuration(elapsed), transferSuffix))
			if progress != nil && progress.OnNodeEnd != nil {
				progress.OnNodeEnd(n.ID, kind, stepErr, elapsed, transferBytes, false)
			}
			return trimLog(b.String()), stepErr
		}
		appendLog(fmt.Sprintf("%s: ok (耗时 %s%s)", linePrefix, formatNodeStepDuration(elapsed), transferSuffix))
		if progress != nil && progress.OnNodeEnd != nil {
			progress.OnNodeEnd(n.ID, kind, nil, elapsed, transferBytes, false)
		}
	}
	return trimLog(b.String()), nil
}

func trimLog(s string) string {
	if len(s) <= maxRunLogBytes {
		return s
	}
	return s[:maxRunLogBytes-20] + "\n...(truncated)"
}

// formatNodeStepDuration 节点步骤耗时，用于日志行「耗时 xxx」
func formatNodeStepDuration(d time.Duration) string {
	if d < time.Second {
		ms := d.Milliseconds()
		if ms < 1 {
			return fmt.Sprintf("%.1fms", float64(d.Microseconds())/1000)
		}
		return fmt.Sprintf("%dms", ms)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return d.Round(time.Second).String()
}

func runShell(ctx context.Context, cwd, script string, procEnv map[string]string, appendLog func(string)) error {
	if strings.TrimSpace(script) == "" {
		return fmt.Errorf("empty script")
	}
	if runtime.GOOS == "windows" {
		// 单条临时 .cmd、一次 cmd /C，与手写 bat 同一会话。对 mvn、npm 等子 .cmd 在各行前自动加
		// call，避免无 call 时父批处理在子 .cmd 处被替换结束。含 &/| 的复合行原样写出，由用户自理。
		bat, err := winBatchFileFromScript(script)
		if err != nil {
			return err
		}
		f, err := os.CreateTemp("", "shh-devflow-*.cmd")
		if err != nil {
			return fmt.Errorf("temp script: %w", err)
		}
		tmpPath := f.Name()
		// 不可在首行前加 UTF-8 BOM：否则首行不是有效的 @echo off，会保持 echo on，把脚本每一行
		// 都回显到输出，看起来像在「逐条执行」；且 BOM 在少数环境下会显示为乱码。
		if _, werr := f.WriteString(bat); werr != nil {
			_ = f.Close()
			_ = os.Remove(tmpPath)
			return fmt.Errorf("write temp script: %w", werr)
		}
		if cerr := f.Close(); cerr != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("close temp script: %w", cerr)
		}
		defer func() { _ = os.Remove(tmpPath) }()
		cmd := exec.CommandContext(ctx, "cmd", "/C", tmpPath)
		cmd.Dir = cwd
		if len(procEnv) > 0 {
			cmd.Env = OsEnvWithProcess(defaultEnviron(), procEnv)
		}
		var outWin cappedMergeWriter
		outWin.limit = maxScriptIOBytes
		cmd.Stdout = &outWin
		cmd.Stderr = &outWin
		err = cmd.Run()
		appendLog(strings.TrimSpace(outWin.resultString()))
		if err != nil {
			return fmt.Errorf("script failed: %w", err)
		}
		return nil
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", script)
	cmd.Dir = cwd
	if len(procEnv) > 0 {
		cmd.Env = OsEnvWithProcess(defaultEnviron(), procEnv)
	}
	var outUnix cappedMergeWriter
	outUnix.limit = maxScriptIOBytes
	cmd.Stdout = &outUnix
	cmd.Stderr = &outUnix
	err := cmd.Run()
	appendLog(strings.TrimSpace(outUnix.resultString()))
	if err != nil {
		return fmt.Errorf("script failed: %w", err)
	}
	return nil
}

func runNodeStep(ctx context.Context, kind string, params map[string]interface{}, nodeID string, nodes []FlowNode, edges []FlowEdge, baseDir string, procEnv map[string]string, appendLog func(string)) (int64, error) {
	switch kind {
	case "git_repo":
		return 0, runGitRepo(ctx, params, baseDir, appendLog)
	case "ssh_connection":
		return 0, runSSHConnection(ctx, params, appendLog)
	case "execute_script":
		cwd := filepath.Join(baseDir, strParam(params, "cwd"))
		if err := os.MkdirAll(cwd, 0o755); err != nil {
			return 0, err
		}
		return 0, runShell(ctx, cwd, strParam(params, "script"), procEnv, appendLog)
	case "remote_ssh_script":
		targets := CollectSshTargetsForRemoteScript(nodeID, nodes, edges)
		if len(targets) == 0 {
			return 0, fmt.Errorf("no upstream SSH (connect SSH or upload node before this step)")
		}
		cwd := strParam(params, "cwd")
		script := strParam(params, "script")
		for i, sp := range targets {
			if len(targets) > 1 {
				appendLog(fmt.Sprintf("--- remote SSH %d/%d: %s@%s ---", i+1, len(targets), strings.TrimSpace(sp.Username), strings.TrimSpace(sp.Host)))
			}
			if err := runRemoteSSHScript(ctx, sp, cwd, script, appendLog); err != nil {
				return 0, err
			}
		}
		return 0, nil
	case "upload_servers":
		sp := sshParamsFromNode(kind, params)
		if sp == nil || strings.TrimSpace(sp.Host) == "" {
			sp = ResolveUpstreamSsh(nodeID, nodes, edges, nil)
		}
		if sp == nil || strings.TrimSpace(sp.Host) == "" {
			return 0, fmt.Errorf("upload: set SSH target on this node or place an upstream SSH node")
		}
		return runUploadServers(ctx, sp, params, baseDir, appendLog)
	case "remote_download":
		sp := ResolveUpstreamSsh(nodeID, nodes, edges, nil)
		if sp == nil || strings.TrimSpace(sp.Host) == "" {
			return 0, fmt.Errorf("no upstream SSH for download")
		}
		return runRemoteDownload(ctx, sp, params, baseDir, appendLog)
	default:
		return 0, fmt.Errorf("unsupported node kind %q", kind)
	}
}

func runGitRepo(ctx context.Context, params map[string]interface{}, baseDir string, appendLog func(string)) error {
	repo := strings.TrimSpace(strParam(params, "repositoryUrl"))
	subdir := strings.TrimSpace(strParam(params, "checkoutSubdir"))
	if subdir == "." {
		subdir = ""
	}
	if repo == "" {
		return fmt.Errorf("repositoryUrl required")
	}
	baseClean := filepath.Clean(baseDir)
	var dest string
	if subdir == "" {
		dest = baseClean
	} else {
		dest = filepath.Join(baseClean, subdir)
		rel, err := filepath.Rel(baseClean, filepath.Clean(dest))
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("checkoutSubdir must stay under workspace")
		}
	}
	_ = os.RemoveAll(dest)
	appendLog("cloning with go-git…")
	err := CloneGitRepoWithGoGit(ctx, params, dest)
	if err != nil {
		return err
	}
	appendLog("clone finished")
	return nil
}

func runSSHConnection(ctx context.Context, params map[string]interface{}, appendLog func(string)) error {
	host := strParam(params, "host")
	port := strParam(params, "port")
	if strings.TrimSpace(port) == "" {
		port = "22"
	}
	user := strParam(params, "username")
	addr := net.JoinHostPort(strings.TrimSpace(host), strings.TrimSpace(port))
	cli, err := DialSSH(ctx, addr, user, strParam(params, "password"), strParam(params, "privateKey"), strParam(params, "authType"))
	if err != nil {
		return err
	}
	_ = cli.Close()
	appendLog("SSH handshake ok")
	return nil
}

func runRemoteSSHScript(ctx context.Context, sp *SshParams, cwd, script string, appendLog func(string)) error {
	if strings.TrimSpace(script) == "" {
		return fmt.Errorf("empty remote script")
	}
	addr := net.JoinHostPort(strings.TrimSpace(sp.Host), strings.TrimSpace(sp.Port))
	cli, err := DialSSH(ctx, addr, sp.Username, sp.Password, sp.PrivateKey, sp.AuthType)
	if err != nil {
		return err
	}
	defer cli.Close()
	sess, err := cli.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	wrap := fmt.Sprintf("set -e\ncd %q\n%s", cwd, script)
	sess.Stdin = strings.NewReader(wrap + "\n")
	var out cappedMergeWriter
	out.limit = maxScriptIOBytes
	sess.Stdout = &out
	sess.Stderr = &out
	err = sess.Run("/bin/bash -s")
	text := strings.TrimSpace(out.resultString())
	if text != "" {
		for _, line := range strings.Split(text, "\n") {
			appendLog(line)
		}
	}
	if err != nil {
		return fmt.Errorf("remote script: %w", err)
	}
	return nil
}

func runUploadServers(ctx context.Context, sp *SshParams, params map[string]interface{}, baseDir string, appendLog func(string)) (int64, error) {
	pattern := strParam(params, "artifactGlob")
	remotePath := strings.TrimRight(strParam(params, "remotePath"), `/`)
	if remotePath == "" {
		return 0, fmt.Errorf("remotePath required")
	}
	matches, err := filepath.Glob(filepath.Join(baseDir, pattern))
	if err != nil {
		return 0, err
	}
	if len(matches) == 0 {
		return 0, fmt.Errorf("no local files match glob %q under %s", pattern, baseDir)
	}
	addr := net.JoinHostPort(strings.TrimSpace(sp.Host), strings.TrimSpace(sp.Port))
	cli, err := DialSSH(ctx, addr, sp.Username, sp.Password, sp.PrivateKey, sp.AuthType)
	if err != nil {
		return 0, err
	}
	defer cli.Close()
	sc, err := sftp.NewClient(cli)
	if err != nil {
		return 0, err
	}
	defer sc.Close()
	var total int64
	for _, local := range matches {
		st, err := os.Stat(local)
		if err != nil {
			return 0, err
		}
		sz := st.Size()
		total += sz
		rel := filepath.Base(local)
		rpath := remotePath + "/" + rel
		if err := sftpPutFile(sc, local, rpath); err != nil {
			return total, fmt.Errorf("upload %s: %w", local, err)
		}
		appendLog(fmt.Sprintf("uploaded %s (%s)", rel, formatByteSize(sz)))
	}
	return total, nil
}

func sftpPutFile(sc *sftp.Client, localPath, remotePath string) error {
	remotePath = filepath.ToSlash(remotePath)
	lf, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer lf.Close()
	_ = sc.MkdirAll(filepath.ToSlash(filepath.Dir(remotePath)))
	rf, err := sc.Create(remotePath)
	if err != nil {
		return err
	}
	defer rf.Close()
	_, err = io.Copy(rf, lf)
	return err
}

func runRemoteDownload(ctx context.Context, sp *SshParams, params map[string]interface{}, baseDir string, appendLog func(string)) (int64, error) {
	remotePath := strParam(params, "remoteSourcePath")
	outRel := strParam(params, "outputPath")
	if remotePath == "" || outRel == "" {
		return 0, fmt.Errorf("remoteSourcePath and outputPath required")
	}
	localPath := filepath.Join(baseDir, outRel)
	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		return 0, err
	}
	addr := net.JoinHostPort(strings.TrimSpace(sp.Host), strings.TrimSpace(sp.Port))
	cli, err := DialSSH(ctx, addr, sp.Username, sp.Password, sp.PrivateKey, sp.AuthType)
	if err != nil {
		return 0, err
	}
	defer cli.Close()
	sc, err := sftp.NewClient(cli)
	if err != nil {
		return 0, err
	}
	defer sc.Close()
	st, err := sc.Stat(remotePath)
	if err != nil {
		return 0, err
	}
	if st.IsDir() {
		return 0, fmt.Errorf("remote path is a directory; download directory not supported in this engine version")
	}
	rf, err := sc.Open(remotePath)
	if err != nil {
		return 0, err
	}
	defer rf.Close()
	lf, err := os.Create(localPath)
	if err != nil {
		return 0, err
	}
	defer lf.Close()
	n, err := io.Copy(lf, rf)
	if err != nil {
		return n, err
	}
	appendLog(fmt.Sprintf("downloaded %s to %s", formatByteSize(n), outRel))
	if boolParam(params, "extract") {
		appendLog("extract flag ignored in engine v1")
	}
	return n, nil
}

// LastExec 状态常量（写入 last_exec_result 列）
const (
	LastExecStatusSuccess    = "success"
	LastExecStatusFailed     = "failed"
	LastExecStatusCancelled  = "cancelled"
)

// 落库最后防线：与 maxRunLogBytes 同量级略放大，防误传超长时仍不撑爆长文本列与前端
const maxLastExecLogBytes = 256 * 1024

// TruncateLastExecLog 截断准备写入 last_exec_log 的文本
func TruncateLastExecLog(s string) string {
	if len(s) <= maxLastExecLogBytes {
		return s
	}
	return s[:maxLastExecLogBytes-24] + "\n...(truncated)"
}

// BuildLastExecRecord 根据 Run 输出与错误生成「状态 + 完整日志文本」供落库
func BuildLastExecRecord(logOut string, runErr error) (status string, logText string) {
	if runErr != nil {
		logOut = logOut + "\nerror: " + runErr.Error()
	}
	logText = TruncateLastExecLog(logOut)
	switch {
	case runErr == nil:
		status = LastExecStatusSuccess
	case errors.Is(runErr, context.Canceled):
		status = LastExecStatusCancelled
	default:
		status = LastExecStatusFailed
	}
	return status, logText
}

// RunStartedAt 用于写 last_exec_time（Unix 秒）
func RunStartedAt() int64 {
	return time.Now().Unix()
}
