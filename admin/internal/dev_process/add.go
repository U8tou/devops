package dev_process

import (
	devimpl "devops/impl/dev_process"
	"devops/model"
	"pkg/constx"
	"pkg/errs"
	r "pkg/resp"
	"pkg/tools/datacv"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

func Add(c *fiber.Ctx) error {
	var rq AddReq
	if err := c.BodyParser(&rq); err != nil {
		return err
	}
	if err := validator.Struct(rq); err != nil {
		return err
	}
	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}
	var args model.DevProcessDto
	args.Code = rq.Code
	args.Remark = rq.Remark
	args.Flow = rq.Flow
	args.CronEnabled = rq.CronEnabled
	args.CronExpr = rq.CronExpr
	if rq.Env != nil {
		j, err := mapToEnvJSON(rq.Env)
		if err != nil {
			return errs.New(err.Error())
		}
		args.EnvJson = j
	}
	args.CreateBy = loginId
	args.UpdateBy = loginId
	args.TagIds = rq.TagIds
	impl := devimpl.Impl()
	affect, err := impl.Add(c.Context(), &args)
	if err != nil {
		return err
	}
	SyncCronForProcess(args.Id)
	return r.Resp(c, AddResp{
		Affect: datacv.IntToStr(affect),
		Id:     datacv.IntToStr(args.Id),
	})
}
