package sysdept

import (
	"pkg/constx"
	r "pkg/resp"
	"pkg/tools/datacv"
	sysdept "system/impl/sys_dept"
	"system/model"
	"time"

	"pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// @Summary		系统部门表_新增
// @Description	新增记录\n权限标识：sys:dept:add
// @Tags			SysDeptApi
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			args	body		AddReq						true	"Request Parameters"
// @Success		200		{object}	r.MyResp{data=AddResp}	"Success"
// @Failure		500		{object}	r.MyResp					"Failure"
// @Router			/sys_dept/add [post]
func Add(c *fiber.Ctx) error {
	// 请求
	var rq AddReq
	err := c.BodyParser(&rq)
	if err != nil {
		return err
	}
	err = validator.Struct(rq)
	if err != nil {
		return err
	}

	loginId, err := constx.GetLoginId(c)
	if err != nil {
		return err
	}
	tp := time.Now().Unix()

	// 处理
	var args model.SysDeptDto
	// args.Id = rq.Id
	args.Pid = datacv.StrToInt(rq.Pid)
	args.Name = rq.Name
	args.Profile = rq.Profile
	args.Leader = rq.Leader
	args.Phone = rq.Phone
	args.Email = rq.Email
	args.Sort = rq.Sort
	args.Status = rq.Status
	args.CreateTime = tp
	args.CreateBy = loginId
	args.UpdateTime = tp
	args.UpdateBy = loginId
	var iSysDeptImpl = sysdept.Impl()
	affect, err := iSysDeptImpl.Add(c.Context(), &args)
	if err != nil {
		return err
	}
	// 响应
	var rp AddResp
	rp.Affect = datacv.IntToStr(affect)
	return r.Resp(c, rp)
}
