package dev_process

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	devimpl "devops/impl/dev_process"
	"pkg/alog"
	"pkg/conf"
	"pkg/devflow"
	"pkg/tools/datacv"

	"github.com/robfig/cron/v3"
)

// 与 pkg/cronvalidate 一致，支持可选「秒」字段的标准 Cron
var devCronParser = cron.NewParser(
	cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

type cronScheduler struct {
	mu      sync.Mutex
	c       *cron.Cron
	entries map[int64]cron.EntryID
	running sync.Map // int64 -> struct{} 防止同一流程并发重叠执行
}

var globalCron *cronScheduler

// StartCronScheduler 在进程内启动全局 Cron（robfig/cron），并从库中加载已启用的定时流程。
func StartCronScheduler() {
	globalCron = &cronScheduler{
		c:       cron.New(cron.WithParser(devCronParser)),
		entries: make(map[int64]cron.EntryID),
	}
	rows, err := devimpl.Impl().ListCronScheduled(context.Background())
	if err != nil {
		alog.Error("dev cron: ListCronScheduled", err)
	} else {
		globalCron.mu.Lock()
		for _, row := range rows {
			globalCron.registerLocked(row.Id, row.CronExpr)
		}
		globalCron.mu.Unlock()
	}
	globalCron.c.Start()
	alog.Info("dev cron: scheduler started", "")
}

// StopCronScheduler 优雅停止 Cron（进程退出前调用）
func StopCronScheduler() {
	if globalCron == nil || globalCron.c == nil {
		return
	}
	ctx := globalCron.c.Stop()
	<-ctx.Done()
	globalCron = nil
	alog.Info("dev cron: scheduler stopped", "")
}

// SyncCronForProcess 在新增/编辑/启用开关后同步单条任务（不符合条件则从调度器移除）
func SyncCronForProcess(id int64) {
	if globalCron == nil {
		return
	}
	row, err := devimpl.Impl().Get(context.Background(), strconv.FormatInt(id, 10))
	globalCron.mu.Lock()
	defer globalCron.mu.Unlock()
	if err != nil || row == nil {
		globalCron.unregisterLocked(id)
		return
	}
	if row.CronEnabled != 1 {
		globalCron.unregisterLocked(id)
		return
	}
	expr := strings.TrimSpace(row.CronExpr)
	if expr == "" {
		globalCron.unregisterLocked(id)
		return
	}
	globalCron.registerLocked(id, expr)
}

// RemoveCronTask 删除流程时从调度器移除（无需再读库）
func RemoveCronTask(id int64) {
	if globalCron == nil {
		return
	}
	globalCron.mu.Lock()
	defer globalCron.mu.Unlock()
	globalCron.unregisterLocked(id)
}

func (s *cronScheduler) unregisterLocked(id int64) {
	if eid, ok := s.entries[id]; ok {
		s.c.Remove(eid)
		delete(s.entries, id)
	}
}

func (s *cronScheduler) registerLocked(id int64, expr string) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return
	}
	s.unregisterLocked(id)
	eid, err := s.c.AddFunc(expr, func() {
		runScheduledProcess(id)
	})
	if err != nil {
		alog.Error("dev cron: AddFunc", fmt.Sprintf("id=%d expr=%q err=%v", id, expr, err))
		return
	}
	s.entries[id] = eid
}

func runScheduledProcess(id int64) {
	if globalCron == nil {
		return
	}
	if _, loaded := globalCron.running.LoadOrStore(id, struct{}{}); loaded {
		alog.Warn("dev cron: skip overlapping run", fmt.Sprintf("id=%d", id))
		return
	}
	go func() {
		defer globalCron.running.Delete(id)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		impl := devimpl.Impl()
		row, err := impl.Get(ctx, strconv.FormatInt(id, 10))
		if err != nil || row == nil {
			return
		}
		if row.CronEnabled != 1 {
			return
		}

		root := strings.TrimSpace(conf.Devops.WorkspaceRoot)
		if root == "" {
			alog.Error("dev cron: devops.workspaceRoot empty", "")
			return
		}

		runStart := time.Now()
		started := runStart.Unix()
		procEnv := envJSONToMap(row.EnvJson)
		idStr := datacv.IntToStr(id)
		logOut, runErr := devflow.Run(ctx, row.Flow, root, idStr, procEnv, nil)
		durationMs := time.Since(runStart).Milliseconds()
		status, logText := devflow.BuildLastExecRecord(logOut, runErr)
		if _, uerr := impl.UpdateLastExec(ctx, id, started, durationMs, status, logText, 0); uerr != nil {
			alog.Error("dev cron: UpdateLastExec", uerr)
		}
		if runErr != nil {
			alog.Error("dev cron: Run", runErr)
		}
	}()
}
