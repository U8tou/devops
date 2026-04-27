package sysuser

import (
	"context"
	"strconv"
	"system/model"
)

// Get 查
func (m *SysUserImpl) Get(ctx context.Context, id string) (*model.SysUserVo, error) {
	var v model.SysUserVo
	_, err := m.engine.Context(ctx).Table(&model.SysUser{}).ID(id).Get(&v)
	if err != nil {
		return nil, err
	}

	// 获取用户关联部门ID
	deptIds, err := m.GetUserDeptIds(ctx, v.Id)
	if err != nil {
		return nil, err
	}
	v.Depts = make([]string, len(deptIds))
	for i, deptId := range deptIds {
		v.Depts[i] = strconv.FormatInt(deptId, 10)
	}

	// 获取用户关联角色ID
	roleIds, err := m.GetUserRoleIds(ctx, v.Id)
	if err != nil {
		return nil, err
	}
	v.Roles = make([]string, len(roleIds))
	for i, roleId := range roleIds {
		v.Roles[i] = strconv.FormatInt(roleId, 10)
	}

	// 获取用户关联岗位ID
	postIds, err := m.GetUserPostIds(ctx, v.Id)
	if err != nil {
		return nil, err
	}
	v.Posts = make([]string, len(postIds))
	for i, postId := range postIds {
		v.Posts[i] = strconv.FormatInt(postId, 10)
	}

	return &v, nil
}

// GetByUserName 按用户名查询（用于登录）
func (m *SysUserImpl) GetByUserName(ctx context.Context, userName string) (*model.SysUserVo, error) {
	var v model.SysUserVo
	has, err := m.engine.Context(ctx).Table(&model.SysUser{}).Where("user_name = ?", userName).Get(&v)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return &v, nil
}
