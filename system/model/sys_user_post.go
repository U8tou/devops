package model

// SysUserPost 用户岗位关联表
type SysUserPost struct {
	UserId int64 `xorm:"pk notnull comment('用户ID')" json:"userId"`
	PostId int64 `xorm:"pk notnull comment('岗位ID')" json:"postId"`
}

func (m *SysUserPost) Comment() string {
	return "用户岗位关联"
}
