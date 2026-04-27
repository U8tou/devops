package sysuser

import (
	"context"
	"fmt"
	"system/model"
)

// Del 软删（xorm deleted 自动设置 delete_time 为秒级时间戳），同时清理关联表
func (m *SysUserImpl) Del(ctx context.Context, ids []int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	for _, id := range ids {
		if id == 1 {
			return 0, fmt.Errorf("不允许删除超级管理员")
		}
	}
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	affect, err := session.In("id", ids).Delete(&model.SysUser{})
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	if _, err = session.In("user_id", ids).Delete(&model.SysUserRole{}); err != nil {
		_ = session.Rollback()
		return 0, err
	}
	if _, err = session.In("user_id", ids).Delete(&model.SysUserDept{}); err != nil {
		_ = session.Rollback()
		return 0, err
	}
	if _, err = session.In("user_id", ids).Delete(&model.SysUserPost{}); err != nil {
		_ = session.Rollback()
		return 0, err
	}

	return affect, session.Commit()
}
