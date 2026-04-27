package sysuser

import (
	"context"
	"system/model"
)

// AssignDept 分配部门
func (m *SysUserImpl) AssignDept(ctx context.Context, userId int64, deptIds []int64) (int64, error) {
	// 开启事务
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 删除旧的用户角色关联
	affect, err := session.Where("user_id = ?", userId).Delete(&model.SysUserDept{})
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	// 新增新的用户角色关联
	if len(deptIds) > 0 {
		userDepts := make([]model.SysUserDept, len(deptIds))
		for i, deptId := range deptIds {
			userDepts[i] = model.SysUserDept{
				UserId: userId,
				DeptId: deptId,
			}
		}
		_, err = session.Insert(&userDepts)
		if err != nil {
			session.Rollback()
			return 0, err
		}
	}

	return affect, session.Commit()
}

// GetUserDeptIds 获取用户部门ID列表
func (m *SysUserImpl) GetUserDeptIds(ctx context.Context, userId int64) ([]int64, error) {
	var deptIds []int64
	err := m.engine.Context(ctx).Table(&model.SysUserDept{}).
		Where("user_id = ?", userId).
		Cols("dept_id").
		Find(&deptIds)
	return deptIds, err
}
