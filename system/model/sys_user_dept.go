package model

// SysUserDept 用户部门关联表
type SysUserDept struct {
	UserId int64 `xorm:"pk notnull comment('用户ID')" json:"userId"`
	DeptId int64 `xorm:"pk notnull comment('部门ID')" json:"deptId"`
}

func (m *SysUserDept) Comment() string {
	return "用户部门关联"
}
