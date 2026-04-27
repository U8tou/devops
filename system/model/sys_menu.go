package model

// SysMenu 系统权限表
type SysMenu struct {
	Id         int64  `xorm:"pk notnull comment('ID')"`
	Pid        int64  `xorm:"notnull index comment('父ID')"`
	Types      int8   `xorm:"default(1) index notnull comment('权限标识:1_菜单,2_权限')"`
	Permis     string `xorm:"varchar(50) comment('权限代码')"`
	Remark     string `xorm:"varchar(50) comment('备注')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
}

func (m *SysMenu) Comment() string {
	return "系统权限"
}

// SysMenuVo 出参
type SysMenuVo struct {
	SysMenu `xorm:"extends"`
	Menus   []int64 `xorm:"extends"`
}

// SysMenuDto 入参
type SysMenuDto struct {
	SysMenu
	Current int `json:"current"`
	Size    int `json:"size"`
	Menus   []int64
}
