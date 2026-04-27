package model

/**
Notes: SysDept MODEL
Time:  2025-04-29 10:48:59
*/

// SysDept 系统部门表
type SysDept struct {
	Id         int64  `xorm:"pk autoincr notnull comment('ID')"`
	Pid        int64  `xorm:"notnull unique(pid_name) comment('上级部门ID')"`
	Name       string `xorm:"varchar(50) notnull unique(pid_name) comment('部门名称')"`
	Profile    string `xorm:"varchar(500) comment('部门简介')"`
	Leader     string `xorm:"varchar(50) notnull comment('负责人')"`
	Phone      string `xorm:"varchar(50) notnull comment('负责人电话')"`
	Email      string `xorm:"varchar(50) notnull comment('负责人邮箱')"`
	Sort       int    `xorm:"default(1) notnull comment('排序')"`
	Status     int8   `xorm:"default(1) comment('状态:1_正常,2_停用')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
	DeleteTime int64  `xorm:"deleted default(0) notnull unique(pid_name) comment('删除时间')"`
}

func (m *SysDept) Comment() string {
	return "系统部门"
}

// SysDeptVo 出参
type SysDeptVo struct {
	SysDept `xorm:"extends"`
}

// SysDeptDto 入参
type SysDeptDto struct {
	SysDept
	Current int `json:"current"`
	Size    int `json:"size"`
}
