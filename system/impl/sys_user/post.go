package sysuser

import (
	"context"
	"system/model"
)

// AssignPost 分配岗位
func (m *SysUserImpl) AssignPost(ctx context.Context, userId int64, postIds []int64) (int64, error) {
	session := m.engine.NewSession()
	defer session.Close()
	session.Context(ctx)
	if err := session.Begin(); err != nil {
		return 0, err
	}

	affect, err := session.Where("user_id = ?", userId).Delete(&model.SysUserPost{})
	if err != nil {
		_ = session.Rollback()
		return 0, err
	}

	if len(postIds) > 0 {
		userPosts := make([]model.SysUserPost, len(postIds))
		for i, postId := range postIds {
			userPosts[i] = model.SysUserPost{
				UserId: userId,
				PostId: postId,
			}
		}
		_, err = session.Insert(&userPosts)
		if err != nil {
			_ = session.Rollback()
			return 0, err
		}
	}

	return affect, session.Commit()
}

// GetUserPostIds 获取用户岗位ID列表
func (m *SysUserImpl) GetUserPostIds(ctx context.Context, userId int64) ([]int64, error) {
	var postIds []int64
	err := m.engine.Context(ctx).Table(&model.SysUserPost{}).
		Where("user_id = ?", userId).
		Cols("post_id").
		Find(&postIds)
	return postIds, err
}
