package dev_project

import (
	devimpl "devops/impl/dev_project"
	r "pkg/resp"
	"pkg/tools/datacv"

	"github.com/gofiber/fiber/v2"
)

func TagList(c *fiber.Ctx) error {
	rows, err := devimpl.Impl().TagList(c.Context())
	if err != nil {
		return err
	}
	out := make([]TagRow, len(rows))
	for i, t := range rows {
		out[i] = TagRow{Id: datacv.IntToStr(t.Id), Name: t.Name}
	}
	return r.Resp(c, TagListResp{Rows: out})
}
