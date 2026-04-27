package dataperm

import (
	"context"
	"pkg/errs"
	sysuser "system/impl/sys_user"
)

// SubjectUserVisible 非根用户的数据范围下，目标用户（subjectUserId）是否可见（与列表 EXISTS 语义一致：目标用户至少有一个部门落在 scope 内）。
func SubjectUserVisible(ctx context.Context, scopeDeptIds []int64, subjectUserId int64) (bool, error) {
	if len(scopeDeptIds) == 0 {
		return false, nil
	}
	deptIds, err := sysuser.Impl().GetUserDeptIds(ctx, subjectUserId)
	if err != nil {
		return false, err
	}
	scopeSet := make(map[int64]struct{}, len(scopeDeptIds))
	for _, id := range scopeDeptIds {
		scopeSet[id] = struct{}{}
	}
	for _, d := range deptIds {
		if _, ok := scopeSet[d]; ok {
			return true, nil
		}
	}
	return false, nil
}

// CheckSubjectUser 非根用户须对目标用户落在数据权限部门内；根用户跳过。
func CheckSubjectUser(ctx context.Context, operatorUserId int64, operatorIsRoot bool, subjectUserId int64) error {
	if operatorIsRoot {
		return nil
	}
	scope, err := ResolveDeptScope(ctx, operatorUserId)
	if err != nil {
		return err
	}
	ok, err := SubjectUserVisible(ctx, scope, subjectUserId)
	if err != nil {
		return err
	}
	if !ok {
		return errs.ERR_SYS_403
	}
	return nil
}

// CheckCreateBy 非根用户须对记录创建者落在数据权限部门内（与列表按 create_by 过滤一致）。
func CheckCreateBy(ctx context.Context, operatorUserId int64, operatorIsRoot bool, createBy int64) error {
	return CheckSubjectUser(ctx, operatorUserId, operatorIsRoot, createBy)
}
