package model

// SysRoleDept 角色部门关联表
type SysRoleDept struct {
	RoleId int64 `xorm:"pk notnull comment('角色ID')" json:"roleId"`
	DeptId int64 `xorm:"pk notnull comment('部门ID')" json:"deptId"`
}

func (m *SysRoleDept) Comment() string {
	return "角色部门关联"
}
