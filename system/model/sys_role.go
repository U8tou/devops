package model

// SysRole 系统角色表
type SysRole struct {
	Id          int64  `xorm:"pk notnull comment('ID')"`
	Name        string `xorm:"varchar(30) index notnull comment('名称')"`
	Role        string `xorm:"varchar(30) unique(role) notnull comment('角色标识')"`
	Status      int8   `xorm:"notnull default(1) comment('状态:1_正常,2_停用')"`
	MenuLinkage int8   `xorm:"notnull comment('操作父子联动:1_联动,2_不联动')"`
	DeptLinkage int8   `xorm:"notnull comment('数据父子联动:1_联动,2_不联动')"`
	Sort        int    `xorm:"notnull default(999) comment('显示顺序')"`
	Remark      string `xorm:"varchar(500) comment('备注')"`
	CreateTime  int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy    int64  `xorm:"notnull comment('创建者')"`
	UpdateTime  int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy    int64  `xorm:"notnull comment('更新者')"`
	DeleteTime  int64  `xorm:"deleted default(0) notnull unique(role) comment('删除时间')"`
}

func (m *SysRole) Comment() string {
	return "系统角色"
}

// SysRoleVo 出参
type SysRoleVo struct {
	SysRole `xorm:"extends"`
	Menus   []int64 `xorm:"extends"`
	Depts   []int64 `xorm:"extends"`
}

// SysRoleDto 入参
type SysRoleDto struct {
	SysRole
	Current int `json:"current"`
	Size    int `json:"size"`
}
