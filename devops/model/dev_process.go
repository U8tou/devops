package model

// DevProcess 自动化流程表
type DevProcess struct {
	Id                 int64  `xorm:"pk autoincr notnull comment('ID')"`
	Code               string `xorm:"varchar(64) notnull unique(dev_process_code) comment('编号')"`
	Remark             string `xorm:"varchar(500) default('') comment('备注')"`
	Flow               string `xorm:"longtext comment('流程JSON')"`
	EnvJson            string `xorm:"longtext notnull default('{}') comment('流程环境变量JSON')"`
	CronExpr           string `xorm:"varchar(128) default('') comment('Cron表达式')"`
	CronEnabled        int8   `xorm:"tinyint notnull default(0) comment('定时启用:0否1是')"`
	LastExecTime       int64  `xorm:"bigint default(0) comment('最近执行时间')"`
	LastExecDurationMs int64  `xorm:"bigint default(0) comment('最近执行总耗时毫秒')"`
	LastExecResult     string `xorm:"varchar(32) default('') comment('最近执行状态:success/failed/cancelled')"`
	LastExecLog        string `xorm:"longtext comment('最近执行日志')"`
	CreateTime         int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy           int64  `xorm:"notnull comment('创建者')"`
	UpdateTime         int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy           int64  `xorm:"notnull comment('更新者')"`
	DeleteTime         int64  `xorm:"deleted default(0) notnull unique(dev_process_code) comment('删除时间')"`
}

func (m *DevProcess) Comment() string {
	return "自动化流程"
}

type DevProcessVo struct {
	DevProcess `xorm:"extends"`
	Tags       []DevProcessTagBrief `xorm:"-" json:"tags"`
}

type DevProcessDto struct {
	DevProcess
	Current int `json:"current"`
	Size    int `json:"size"`
	// TagIds 新增/编辑时提交的合法标签 id（仅存在于字典中的）
	TagIds []int64 `xorm:"-" json:"-"`
	// TagFilterIds、TagFilterOther 列表筛选（OR 语义）
	TagFilterIds   []int64 `xorm:"-" json:"-"`
	TagFilterOther bool    `xorm:"-" json:"-"`
	// TagFilterExclude 为 true 时表示「不包含」所选标签条件（对整条标签 OR 条件取反）
	TagFilterExclude bool `xorm:"-" json:"-"`
	DataScopeActive  bool    `xorm:"-" json:"-"`
	DataScopeDeptIds []int64 `xorm:"-" json:"-"`
}

// DevCronScheduleRow 仅用于列出已启用的定时流程（id + cron_expr）
type DevCronScheduleRow struct {
	Id       int64  `xorm:"id"`
	CronExpr string `xorm:"cron_expr"`
}
