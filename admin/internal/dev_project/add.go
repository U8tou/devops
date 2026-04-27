package dev_project

import (
	devimpl "devops/impl/dev_project"
	"devops/model"
	"pkg/constx"
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
	var args model.DevProjectDto
	args.Name = rq.Name
	args.Status = rq.Status
	args.Progress = rq.Progress
	args.VersionChangelog = rq.VersionChangelog
	args.MindJson = rq.MindJson
	args.TagIds = rq.TagIds
	args.CreateBy = loginId
	args.UpdateBy = loginId
	impl := devimpl.Impl()
	affect, err := impl.Add(c.Context(), &args)
	if err != nil {
		return err
	}
	return r.Resp(c, AddResp{
		Affect: datacv.IntToStr(affect),
		Id:     datacv.IntToStr(args.Id),
	})
}
