package dev_process

import (
	"context"
	"strings"
	"time"

	"admin/internal/datapermctx"
	devimpl "devops/impl/dev_process"
	"pkg/conf"
	"pkg/devflow"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type RunReq struct {
	Id string `json:"id" validate:"required,number"`
}

type RunResp struct {
	Ok  bool   `json:"ok"`
	Log string `json:"log"`
}

func Run(c *fiber.Ctx) error {
	var rq RunReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
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
	row, err := impl.Get(c.Context(), rq.Id)
	if err != nil {
		return err
	}
	if row == nil {
		return errs.ERR_DB_NO_EXIST
	}
	if err := datapermctx.CheckCreateByWith(c, loginId, isRoot, row.CreateBy); err != nil {
		return err
	}

	ctx := c.UserContext()
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()

	runStart := time.Now()
	started := runStart.Unix()
	procEnv := envJSONToMap(row.EnvJson)
	logOut, runErr := devflow.Run(ctx, row.Flow, root, rq.Id, procEnv, nil)
	durationMs := time.Since(runStart).Milliseconds()
	status, logText := devflow.BuildLastExecRecord(logOut, runErr)
	pid := datacv.StrToInt(rq.Id)
	if _, uerr := impl.UpdateLastExec(c.Context(), pid, started, durationMs, status, logText, loginId); uerr != nil {
		return uerr
	}
	if runErr != nil {
		return errs.Sys(runErr)
	}
	return r.Resp(c, RunResp{Ok: true, Log: logOut})
}
