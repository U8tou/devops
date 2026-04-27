package datapermctx

import (
	devopsmodel "devops/model"
	"pkg/auth"
	"pkg/constx"
	"pkg/tools/datacv"
	dataperm "system/impl/dataperm"
	"system/model"

	"github.com/gofiber/fiber/v2"
)

// LoginRoot 当前请求登录用户 ID 与是否根用户。
func LoginRoot(c *fiber.Ctx) (loginId int64, isRoot bool, err error) {
	loginId, err = constx.GetLoginId(c)
	if err != nil {
		return 0, false, err
	}
	isRoot = auth.IsRootUser(c.Context(), datacv.IntToStr(loginId))
	return loginId, isRoot, nil
}

func applyScope(c *fiber.Ctx, fn func(active bool, ids []int64)) error {
	loginId, isRoot, err := LoginRoot(c)
	if err != nil {
		return err
	}
	if isRoot {
		return nil
	}
	scope, err := dataperm.ResolveDeptScope(c.Context(), loginId)
	if err != nil {
		return err
	}
	fn(true, scope)
	return nil
}

// ApplyUserDto 列表/导出等与用户相关的查询 DTO 注入数据范围。
func ApplyUserDto(c *fiber.Ctx, args *model.SysUserDto) error {
	return applyScope(c, func(active bool, ids []int64) {
		args.DataScopeActive = active
		args.DataScopeDeptIds = ids
	})
}

// ApplyDevProjectDto 项目列表等查询 DTO 注入数据范围。
func ApplyDevProjectDto(c *fiber.Ctx, args *devopsmodel.DevProjectDto) error {
	return applyScope(c, func(active bool, ids []int64) {
		args.DataScopeActive = active
		args.DataScopeDeptIds = ids
	})
}

// ApplyDevProcessDto 流程列表等查询 DTO 注入数据范围。
func ApplyDevProcessDto(c *fiber.Ctx, args *devopsmodel.DevProcessDto) error {
	return applyScope(c, func(active bool, ids []int64) {
		args.DataScopeActive = active
		args.DataScopeDeptIds = ids
	})
}

// CheckSubjectUser 校验目标系统用户是否在当前用户数据范围内。
func CheckSubjectUser(c *fiber.Ctx, subjectUserId int64) error {
	loginId, isRoot, err := LoginRoot(c)
	if err != nil {
		return err
	}
	return CheckSubjectUserWith(c, loginId, isRoot, subjectUserId)
}

// CheckSubjectUserWith 使用已解析的 loginId / isRoot，避免重复读取会话。
func CheckSubjectUserWith(c *fiber.Ctx, loginId int64, isRoot bool, subjectUserId int64) error {
	return dataperm.CheckSubjectUser(c.Context(), loginId, isRoot, subjectUserId)
}

// CheckCreateBy 校验业务记录创建者是否在当前用户数据范围内。
func CheckCreateBy(c *fiber.Ctx, createBy int64) error {
	loginId, isRoot, err := LoginRoot(c)
	if err != nil {
		return err
	}
	return CheckCreateByWith(c, loginId, isRoot, createBy)
}

// CheckCreateByWith 使用已解析的 loginId / isRoot。
func CheckCreateByWith(c *fiber.Ctx, loginId int64, isRoot bool, createBy int64) error {
	return dataperm.CheckCreateBy(c.Context(), loginId, isRoot, createBy)
}
