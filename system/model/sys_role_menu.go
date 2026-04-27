package model

// SysRoleMenu 角色菜单关联表
type SysRoleMenu struct {
	RoleId int64 `xorm:"pk notnull comment('角色ID')" json:"roleId"`
	MenuId int64 `xorm:"pk notnull comment('菜单ID')" json:"menuId"`
}

func (m *SysRoleMenu) Comment() string {
	return "角色菜单关联"
}
