package sysuser

import (
	"context"
	"pkg/db"
	"sync"
	"system/model"

	"xorm.io/xorm"
)

/**
Notes: 系统用户 IMPL
Time:  2025-12-19 11:23:48
*/

var (
	sysUserImpl     *SysUserImpl
	sysUserImplOnce sync.Once
)

// SysUserImpl 类
type SysUserImpl struct {
	engine *xorm.Engine
}

// Impl 实例化 | 依赖注入
func Impl() ISysUserImpl {
	sysUserImplOnce.Do(func() {
		sysUserImpl = &SysUserImpl{
			engine: db.GetDb(),
		}
	})
	return sysUserImpl
}

// ISysUserImpl 接口
type ISysUserImpl interface {
	// List 分页
	List(ctx context.Context, args *model.SysUserDto) (int64, []model.SysUserVo, error)
	// Get 查询
	Get(ctx context.Context, id string) (*model.SysUserVo, error)
	// GetByUserName 按用户名查询（用于登录）
	GetByUserName(ctx context.Context, userName string) (*model.SysUserVo, error)
	// Del 删除
	Del(ctx context.Context, ids []int64) (int64, error)
	// Add 新增
	Add(ctx context.Context, args *model.SysUserDto) (int64, error)
	// Edit 编辑
	Edit(ctx context.Context, args *model.SysUserDto) (int64, error)
	// ResetPassword 重置密码（管理员）
	ResetPassword(ctx context.Context, userId int64, newPassword string) (int64, error)
	// AssignRole 分配角色
	AssignRole(ctx context.Context, userId int64, roleIds []int64) (int64, error)
	// GetUserRoleIds 获取用户角色ID列表
	GetUserRoleIds(ctx context.Context, userId int64) ([]int64, error)
	// AssignDept 分配部门
	AssignDept(ctx context.Context, userId int64, roleIds []int64) (int64, error)
	// GetUserDeptIds 获取用户部门ID列表
	GetUserDeptIds(ctx context.Context, userId int64) ([]int64, error)
	// AssignPost 分配岗位
	AssignPost(ctx context.Context, userId int64, postIds []int64) (int64, error)
	// GetUserPostIds 获取用户岗位ID列表
	GetUserPostIds(ctx context.Context, userId int64) ([]int64, error)
}
