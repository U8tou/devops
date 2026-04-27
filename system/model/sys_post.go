package model

/**
Notes: SysPost MODEL
Time:  2025-04-29 10:48:59
*/

// SysPost 系统岗位表
type SysPost struct {
	Id         int64  `xorm:"pk autoincr notnull comment('ID')"`
	Name       string `xorm:"varchar(50) notnull unique(post_name) comment('岗位名称')"`
	Sort       int    `xorm:"default(999) notnull comment('排序')"`
	Status     int8   `xorm:"default(1) comment('状态:1_正常,2_停用')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
	DeleteTime int64  `xorm:"deleted default(0) notnull unique(post_name) comment('删除时间')"`
}

func (m *SysPost) Comment() string {
	return "系统岗位"
}

// SysPostVo 出参
type SysPostVo struct {
	SysPost `xorm:"extends"`
}

// SysPostDto 入参
type SysPostDto struct {
	SysPost
	Current int `json:"current"`
	Size    int `json:"size"`
}
