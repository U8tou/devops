package sysrole

import (
	"context"
	"pkg/db"
	"sync"
	"system/model"

	"xorm.io/xorm"
)

/**
Notes: 系统角色 IMPL
Time:  2025-01-29
*/

var (
	sysRoleImpl     *SysRoleImpl
	sysRoleImplOnce sync.Once
)

// SysRoleImpl 类
type SysRoleImpl struct {
	engine *xorm.Engine
}

// Impl 实例化 | 依赖注入
func Impl() ISysRoleImpl {
	sysRoleImplOnce.Do(func() {
		sysRoleImpl = &SysRoleImpl{
			engine: db.GetDb(),
		}
	})
	return sysRoleImpl
}

// ISysRoleImpl 接口
type ISysRoleImpl interface {
	// List 分页
	List(ctx context.Context, args *model.SysRoleDto) (int64, []model.SysRoleVo, error)
	// Get 查询
	Get(ctx context.Context, id int64) (*model.SysRoleVo, error)
	// Del 删除
	Del(ctx context.Context, ids []int64) (int64, error)
	// Add 新增
	Add(ctx context.Context, args *model.SysRoleDto) (int64, error)
	// Edit 编辑
	Edit(ctx context.Context, args *model.SysRoleDto) (int64, error)
	// All 全量
	All(ctx context.Context) ([]model.SysRoleVo, error)
	// AssignDept 分配部门权限
	AssignDept(ctx context.Context, deptId int64, deptIds []int64, deptLinkage int8) (int64, error)
	// AssignMenu 分配菜单权限
	AssignMenu(ctx context.Context, roleId int64, menuIds []int64, menuLinkage int8) (int64, error)
}
