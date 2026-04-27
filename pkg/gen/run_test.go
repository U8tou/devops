package gen

import (
	"pkg/gen/core"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestSomething(t *testing.T) {
	// 代码生成: 1.待解析结构体 2.表描述(中文,如系统岗位)
	// 结构体命名规则: 模块+表名,
	// 如 SysUser(system简写Sys), 模块名就为system
	// 如 AppUser, 模块名就为app
	// 导入包时，sys/** app/**
	core.CreateCode(SysPost{}, "系统岗位")
}

type SysUser struct {
	Id         int64  `xorm:"pk autoincr comment('用户ID')"`
	DeptId     int64  `xorm:"index comment('部门ID')"`
	UserName   string `xorm:"varchar(30) notnull unique('user_name') comment('用户账号')"`
	NickName   string `xorm:"varchar(30) notnull index comment('用户昵称')"`
	UserType   string `xorm:"varchar(2) default('00') notnull comment('用户类型:00_系统用户')"`
	Email      string `xorm:"varchar(50) notnull index comment('用户邮箱')"`
	PhoneArea  string `xorm:"varchar(10) notnull default('+86') comment('电话区号')"`
	Phone      string `xorm:"varchar(20) notnull index comment('电话号码')"`
	Sex        int8   `xorm:"default(1) notnull index comment('用户性别:1_男,2_女,3_未知')"`
	Avatar     string `xorm:"varchar(100) notnull comment('头像')"`
	Password   string `xorm:"varchar(100) notnull comment('密码')"`
	Status     int8   `xorm:"varchar(1) notnull index default(1) comment('状态:1_正常,2_停用')"`
	Address    string `xorm:"varchar(255) default('') notnull comment('联系地址')"`
	Remark     string `xorm:"varchar(500) default('') notnull comment('备注')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
	DeleteTime int64  `xorm:"unique('user_name') default(0) notnull comment('删除标记:0_未删除,非0_已删除(存ID保证唯一)')"`
}

type SysDept struct {
	Id         int64  `xorm:"pk autoincr notnull comment('ID')"`
	Pid        int64  `xorm:"notnull unique('pid_name') comment('上级部门ID')"`
	Name       string `xorm:"varchar(50) notnull unique('pid_name') comment('部门名称')"`
	Profile    string `xorm:"varchar(500) comment('部门简介')"`
	Leader     string `xorm:"varchar(50) notnull comment('负责人')"`
	Phone      string `xorm:"varchar(50) index notnull comment('负责人电话')"`
	Email      string `xorm:"varchar(50) index notnull comment('负责人邮箱')"`
	Sort       int    `xorm:"default(999) notnull comment('排序')"`
	Status     int8   `xorm:"default(1) index comment('状态:1_正常,2_停用')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
	DeleteTime int64  `xorm:"unique('pid_name') default(0) notnull comment('删除标记:0_未删除,非0_已删除(存ID保证唯一)')"`
}

// SysPost 系统岗位表
type SysPost struct {
	Id         int64  `xorm:"pk autoincr notnull comment('ID')"`
	Name       string `xorm:"varchar(50) notnull unique('post_name') comment('岗位名称')"`
	DeleteTime int64  `xorm:"unique('post_name') default(0) notnull comment('删除标记:0_未删除,非0_已删除(存ID保证唯一)')"`
	Sort       int    `xorm:"default(999) notnull comment('排序')"`
	Status     int8   `xorm:"default(1) comment('状态:1_正常,2_停用')"`
	CreateTime int64  `xorm:"created notnull comment('创建时间')"`
	CreateBy   int64  `xorm:"notnull comment('创建者')"`
	UpdateTime int64  `xorm:"updated notnull comment('更新时间')"`
	UpdateBy   int64  `xorm:"notnull comment('更新者')"`
}
