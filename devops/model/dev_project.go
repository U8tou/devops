package model

// DevProject 项目管理表
type DevProject struct {
	Id               int64  `xorm:"pk autoincr notnull comment('ID')"`
	Name             string `xorm:"varchar(200) notnull unique(dev_project_name) comment('项目名称')"`
	Status           int8   `xorm:"tinyint notnull default(0) comment('状态:0草稿1进行中2暂停3已完成')"`
	Progress         int8   `xorm:"tinyint notnull default(0) comment('整体进度0-100')"`
	VersionChangelog string `xorm:"longtext comment('版本更新日志')"`
	MindJson         string `xorm:"longtext notnull default('{}') comment('需求思维导图JSON')"`
	CreateTime       int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy         int64  `xorm:"notnull comment('创建者')"`
	UpdateTime       int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy         int64  `xorm:"notnull comment('更新者')"`
	DeleteTime       int64  `xorm:"deleted default(0) notnull unique(dev_project_name) comment('删除时间')"`
}

func (m *DevProject) Comment() string {
	return "项目管理"
}

// DevProjectVo 列表/详情（含标签）
type DevProjectVo struct {
	DevProject `xorm:"extends"`
	Tags       []DevProjectTagBrief `xorm:"-" json:"tags"`
}

// DevProjectDto 列表/查询参数
type DevProjectDto struct {
	DevProject
	Current      int    `json:"current"`
	Size         int    `json:"size"`
	FilterName   string `xorm:"-" json:"-"` // 筛选：名称模糊
	FilterStatus *int8  `xorm:"-" json:"-"` // 筛选：状态，nil 表示不限
	// TagIds 新增/编辑时提交的合法标签 id
	TagIds []int64 `xorm:"-" json:"-"`
	// TagFilterIds、TagFilterOther 列表筛选（OR 语义）
	TagFilterIds     []int64 `xorm:"-" json:"-"`
	TagFilterOther   bool    `xorm:"-" json:"-"`
	TagFilterExclude bool    `xorm:"-" json:"-"`
	DataScopeActive  bool    `xorm:"-" json:"-"`
	DataScopeDeptIds []int64 `xorm:"-" json:"-"`
}
