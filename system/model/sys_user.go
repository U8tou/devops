package model

/**
Notes: 系统用户 MODEL
Time:  2025-12-18 17:46:03
*/

// 用户类型（与 sys_user.user_type 一致）
const (
	UserTypeRoot   = "00" // 根用户：全量菜单权限，不参与列表
	UserTypeNormal = "10" // 普通用户：需绑定角色后方可登录
)

// SysUser 系统用户
type SysUser struct {
	Id         int64  `xorm:"pk autoincr comment('用户ID')"`
	UserName   string `xorm:"varchar(30) notnull unique(user_name) comment('用户账号')"`
	NickName   string `xorm:"varchar(30) notnull index comment('用户昵称')"`
	UserType   string `xorm:"varchar(2) default('00') notnull comment('用户类型:00_根管理员,10_普通账户')"`
	Email      string `xorm:"varchar(50) notnull index comment('用户邮箱')"`
	PhoneArea  string `xorm:"varchar(10) notnull default('+86') comment('电话区号')"`
	Phone      string `xorm:"varchar(20) notnull index comment('电话号码')"`
	Sex        int8   `xorm:"default(1) notnull index comment('用户性别:1_男,2_女,3_未知')"`
	Avatar     string `xorm:"varchar(100) notnull comment('头像')"`
	Password   string `xorm:"varchar(100) notnull comment('密码')"`
	Status     int8   `xorm:"notnull index default(1) comment('状态:1_正常,2_停用')"`
	Address    string `xorm:"varchar(255) default('') notnull comment('联系地址')"`
	Remark     string `xorm:"varchar(500) default('') notnull comment('备注')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
	DeleteTime int64  `xorm:"deleted default(0) notnull unique(user_name) comment('删除时间')"`
}

func (m *SysUser) Comment() string {
	return "系统用户"
}

// SysUserVo 出参
type SysUserVo struct {
	SysUser `xorm:"extends"`
	Depts   []string `xorm:"extends"`
	Roles   []string `xorm:"extends"`
	Posts   []string `xorm:"extends"`
}

// SysUserDto 入参
type SysUserDto struct {
	SysUser
	Current         int   `json:"current"`
	Size            int   `json:"size"`
	DeptId          int64 `json:"deptId"`          // 部门ID（用于查询筛选）
	CreateTimeStart int64 `json:"createTimeStart"` // 创建时间开始（秒级时间戳）
	CreateTimeEnd   int64 `json:"createTimeEnd"`   // 创建时间结束（秒级时间戳）
	// DataScopeActive 为 true 时按 DataScopeDeptIds 过滤（非根用户）；空切片表示无数据权限
	DataScopeActive  bool    `json:"-"`
	DataScopeDeptIds []int64 `json:"-"`
}
