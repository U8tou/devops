package model

// SysUserRole 用户角色关联表
type SysUserRole struct {
	UserId int64 `xorm:"pk notnull comment('用户ID')" json:"userId"`
	RoleId int64 `xorm:"pk notnull comment('角色ID')" json:"roleId"`
}

func (m *SysUserRole) Comment() string {
	return "用户角色关联"
}
