package dev_process

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_process"
	"pkg/conf"
	"pkg/devflow"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

// runCancelMap：同一用户 + 流程 id 仅保留一次执行；新开流会先取消旧 run。
var runCancelMap sync.Map // key: runCancelKey -> context.CancelFunc

func runCancelKey(loginId int64, processID string) string {
	return fmt.Sprintf("%d:%s", loginId, processID)
}

type sseProgressPayload struct {
	Type          string `json:"type"`
	NodeID        string `json:"nodeId,omitempty"`
	Kind          string `json:"kind,omitempty"`
	Line          string `json:"line,omitempty"`
	OK            bool   `json:"ok,omitempty"`
	Error         string `json:"error,omitempty"`
	StartedAtMs   int64  `json:"startedAtMs,omitempty"`   // 节点开始执行（Unix 毫秒），node_start
	DurationMs    int64  `json:"durationMs,omitempty"`    // 节点耗时（毫秒），node_end
	TransferBytes int64  `json:"transferBytes,omitempty"` // 上传/下载字节数，node_end
	Skipped       bool   `json:"skipped,omitempty"`       // 本节点因流程开关未执行
}

type sseDonePayload struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func appendSSE(event string, data []byte) []byte {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "event: %s\ndata: %s\n\n", event, data)
	return buf.Bytes()
}

// RunStream GET /dev_process/run_stream?id=流程ID — SSE 实时推送执行进度与日志（Query 可带 token 供 EventSource）。
func RunStream(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Query("id"))
	if id == "" {
		return errs.New("id required")
	}
	root := strings.TrimSpace(conf.Devops.WorkspaceRoot)
	if root == "" {
		return errs.New("devops.workspaceRoot 未配置，请在 setting.yaml 中设置 devops.workspaceRoot")
	}
	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	impl := devimpl.Impl()
	row, err := impl.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	if err := datapermctx.CheckCreateByWith(c, loginId, isRoot, row.CreateBy); err != nil {
		return err
	}

	flow := row.Flow
	procEnv := envJSONToMap(row.EnvJson)
	pid := datacv.StrToInt(id)
	userCtx := c.UserContext()
	if userCtx == nil {
		userCtx = context.Background()
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		ctx, cancel := context.WithTimeout(userCtx, 30*time.Minute)
		key := runCancelKey(loginId, id)
		if old, ok := runCancelMap.LoadAndDelete(key); ok {
			if fn, ok := old.(context.CancelFunc); ok {
				fn()
			}
		}
		runCancelMap.Store(key, cancel)

		events := make(chan []byte, 256)

		go func() {
			defer close(events)
			defer func() {
				runCancelMap.Delete(key)
				cancel()
			}()

			send := func(p sseProgressPayload) {
				b, err := json.Marshal(p)
				if err != nil {
					return
				}
				var buf bytes.Buffer
				fmt.Fprintf(&buf, "event: progress\ndata: %s\n\n", b)
				select {
				case events <- buf.Bytes():
				case <-ctx.Done():
				}
			}

			progress := &devflow.RunProgress{
				OnLog: func(line string) {
					send(sseProgressPayload{Type: "log", Line: line})
				},
				OnNodeStart: func(nodeID, kind string, startedAt time.Time) {
					send(sseProgressPayload{
						Type:        "node_start",
						NodeID:      nodeID,
						Kind:        kind,
						StartedAtMs: startedAt.UnixMilli(),
					})
				},
				OnNodeEnd: func(nodeID, kind string, stepErr error, duration time.Duration, transferBytes int64, skipped bool) {
					p := sseProgressPayload{
						Type:          "node_end",
						NodeID:        nodeID,
						Kind:          kind,
						OK:            stepErr == nil,
						DurationMs:    duration.Milliseconds(),
						TransferBytes: transferBytes,
						Skipped:       skipped,
					}
					if stepErr != nil {
						p.Error = stepErr.Error()
					}
					send(p)
				},
			}

			runStart := time.Now()
			started := runStart.Unix()
			logOut, runErr := devflow.Run(ctx, flow, root, id, procEnv, progress)
			durationMs := time.Since(runStart).Milliseconds()
			status, logText := devflow.BuildLastExecRecord(logOut, runErr)
			// StreamWriter 在独立 goroutine 中，勿使用已失效的 fiber RequestCtx
			if _, uerr := impl.UpdateLastExec(context.Background(), pid, started, durationMs, status, logText, loginId); uerr != nil {
				d, _ := json.Marshal(sseDonePayload{OK: false, Error: uerr.Error()})
				events <- appendSSE("run_error", d)
				return
			}

			if runErr != nil {
				if errors.Is(runErr, context.Canceled) {
					d, _ := json.Marshal(map[string]string{"reason": "cancelled"})
					events <- appendSSE("cancelled", d)
					return
				}
				d, _ := json.Marshal(sseDonePayload{OK: false, Error: runErr.Error()})
				events <- appendSSE("done", d)
				return
			}
			d, _ := json.Marshal(sseDonePayload{OK: true})
			events <- appendSSE("done", d)
		}()

		for chunk := range events {
			if _, err := w.Write(chunk); err != nil {
				return
			}
			if err := w.Flush(); err != nil {
				return
			}
		}
	}))
	return nil
}

// RunCancel POST /dev_process/run_cancel?id= — 取消当前用户对该流程的流式执行（与 RunStream 共用 cancel）。
func RunCancel(c *fiber.Ctx) error {
	id := strings.TrimSpace(c.Query("id"))
	if id == "" {
		return errs.New("id required")
	}
	loginId, isRoot, err := datapermctx.LoginRoot(c)
	if err != nil {
		return err
	}
	impl := devimpl.Impl()
	row, err := impl.Get(c.Context(), id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	if err := datapermctx.CheckCreateByWith(c, loginId, isRoot, row.CreateBy); err != nil {
		return err
	}
	key := runCancelKey(loginId, id)
	v, loaded := runCancelMap.LoadAndDelete(key)
	if loaded {
		if fn, ok := v.(context.CancelFunc); ok {
			fn()
		}
	}
	return r.Resp(c, map[string]bool{"ok": true})
}
