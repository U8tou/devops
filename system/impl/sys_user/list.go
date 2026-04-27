package sysuser

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"system/model"
)

// List 分页
func (m *SysUserImpl) List(ctx context.Context, args *model.SysUserDto) (int64, []model.SysUserVo, error) {
	// 构建筛选条件（排除软删除）；根用户不在列表中展示
	query := m.engine.Context(ctx).Table(&model.SysUser{}).Where("user_type <> ?", model.UserTypeRoot)
	if args.DataScopeActive {
		if len(args.DataScopeDeptIds) == 0 {
			query.And("1 = 0")
		} else {
			placeholders := strings.Repeat("?,", len(args.DataScopeDeptIds))
			placeholders = placeholders[:len(placeholders)-1]
			scopeBind := make([]any, len(args.DataScopeDeptIds))
			for i, id := range args.DataScopeDeptIds {
				scopeBind[i] = id
			}
			query.And(fmt.Sprintf("EXISTS (SELECT 1 FROM sys_user_dept u WHERE u.user_id = sys_user.id AND u.dept_id IN (%s))", placeholders), scopeBind...)
		}
	}
	if args.Id != 0 {
		query.And("id = ?", args.Id)
	}
	if args.UserName != "" {
		query.And("username LIKE ?", "%"+args.UserName+"%")
	}
	if args.NickName != "" {
		query.And("nickname LIKE ?", "%"+args.NickName+"%")
	}
	if args.Email != "" {
		query.And("email LIKE ?", "%"+args.Email+"%")
	}
	if args.Phone != "" {
		query.And("phone LIKE ?", "%"+args.Phone+"%")
	}
	if args.Sex != 0 {
		query.And("sex = ?", args.Sex)
	}
	if args.Status != 0 {
		query.And("status = ?", args.Status)
	}
	if args.CreateTimeStart > 0 {
		query.And("create_time >= ?", args.CreateTimeStart)
	}
	if args.CreateTimeEnd > 0 {
		query.And("create_time <= ?", args.CreateTimeEnd)
	}
	if args.DeptId != 0 {
		// 获取指定部门及其所有子部门的ID列表
		deptIds, err := m.getDeptAndChildrenIds(ctx, args.DeptId)
		if err != nil {
			return 0, nil, err
		}
		if len(deptIds) > 0 {
			// 通过 EXISTS 子查询来筛选部门ID（包含该部门及其所有子部门），避免 JOIN 导致的重复记录
			// 构建 IN 查询的占位符和参数
			placeholders := strings.Repeat("?,", len(deptIds))
			placeholders = placeholders[:len(placeholders)-1] // 移除最后一个逗号
			args := make([]any, len(deptIds))
			for i, id := range deptIds {
				args[i] = id
			}
			query.And(fmt.Sprintf("EXISTS (SELECT 1 FROM sys_user_dept WHERE sys_user_dept.user_id = sys_user.id AND sys_user_dept.dept_id IN (%s))", placeholders), args...)
		}
	}

	// 查询
	var total int64
	var err error
	rows := make([]model.SysUserVo, 0)
	if args.Size > 0 {
		total, err = query.Limit(args.Size, (args.Current-1)*args.Size).Asc("id").FindAndCount(&rows)
	} else {
		err = query.Asc("id").Find(&rows)
	}
	if err != nil {
		return 0, nil, err
	}

	// 为每个用户获取关联的部门ID列表和角色ID列表
	for i := range rows {
		// 获取用户关联部门ID
		deptIds, err := m.GetUserDeptIds(ctx, rows[i].Id)
		if err != nil {
			return 0, nil, err
		}
		rows[i].Depts = make([]string, len(deptIds))
		for j, deptId := range deptIds {
			rows[i].Depts[j] = strconv.FormatInt(deptId, 10)
		}

		// 获取用户关联角色ID
		roleIds, err := m.GetUserRoleIds(ctx, rows[i].Id)
		if err != nil {
			return 0, nil, err
		}
		rows[i].Roles = make([]string, len(roleIds))
		for j, roleId := range roleIds {
			rows[i].Roles[j] = strconv.FormatInt(roleId, 10)
		}

		// 获取用户关联岗位ID
		postIds, err := m.GetUserPostIds(ctx, rows[i].Id)
		if err != nil {
			return 0, nil, err
		}
		rows[i].Posts = make([]string, len(postIds))
		for j, postId := range postIds {
			rows[i].Posts[j] = strconv.FormatInt(postId, 10)
		}
	}

	return total, rows, nil
}

// getDeptAndChildrenIds 递归获取指定部门及其所有子部门的ID列表
func (m *SysUserImpl) getDeptAndChildrenIds(ctx context.Context, deptId int64) ([]int64, error) {
	var result []int64
	result = append(result, deptId) // 包含自身

	// 递归获取所有子部门
	var children []model.SysDept
	err := m.engine.Context(ctx).Table(&model.SysDept{}).
		Where("pid = ?", deptId).
		Find(&children)
	if err != nil {
		return nil, err
	}

	// 递归处理每个子部门
	for _, child := range children {
		childIds, err := m.getDeptAndChildrenIds(ctx, child.Id)
		if err != nil {
			return nil, err
		}
		result = append(result, childIds...)
	}

	return result, nil
}
