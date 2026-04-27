package dataperm

import (
	"context"
	"pkg/db"
	"system/model"
	sysrole "system/impl/sys_role"
	sysuser "system/impl/sys_user"

	"xorm.io/xorm"
)

const (
	// DeptLinkageWithChildren 数据权限含子部门（与 model.SysRole 注释一致）
	DeptLinkageWithChildren int8 = 1
)

// ResolveDeptScope 解析登录用户通过「角色-数据权限部门」可访问的部门 ID（多角色并集、按各角色 dept_linkage 展开子部门、去重）。
// 用于非根用户；若角色均未配置部门，返回空切片。
func ResolveDeptScope(ctx context.Context, userId int64) ([]int64, error) {
	roleIds, err := sysuser.Impl().GetUserRoleIds(ctx, userId)
	if err != nil {
		return nil, err
	}
	if len(roleIds) == 0 {
		return nil, nil
	}
	engine := db.GetDb()
	seen := make(map[int64]struct{})
	for _, rid := range roleIds {
		roleVo, err := sysrole.Impl().Get(ctx, rid)
		if err != nil {
			return nil, err
		}
		if roleVo == nil {
			continue
		}
		deptIds := roleVo.Depts
		for _, did := range deptIds {
			if roleVo.DeptLinkage == DeptLinkageWithChildren {
				expanded, err := getDeptAndChildrenIds(ctx, engine, did)
				if err != nil {
					return nil, err
				}
				for _, eid := range expanded {
					seen[eid] = struct{}{}
				}
			} else {
				seen[did] = struct{}{}
			}
		}
	}
	out := make([]int64, 0, len(seen))
	for id := range seen {
		out = append(out, id)
	}
	return out, nil
}

func getDeptAndChildrenIds(ctx context.Context, engine *xorm.Engine, deptId int64) ([]int64, error) {
	result := []int64{deptId}
	var children []model.SysDept
	err := engine.Context(ctx).Table(&model.SysDept{}).Where("pid = ?", deptId).Find(&children)
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		sub, err := getDeptAndChildrenIds(ctx, engine, child.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, sub...)
	}
	return result, nil
}
