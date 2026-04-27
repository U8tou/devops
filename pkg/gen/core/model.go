package core

import "strings"

type args struct {
	StructName      string      // 结构体名
	StructDoc       string      // 结构体描述(表名/中文描述)
	ModulePath      string      // 模块路径(下划线命名,如 sys_user)
	EntityName      string      // 实体名(Pascal,如 Post)
	EntityKebab     string      // 实体名(kebab,如 post)
	ViewParent      string      // 视图父路径(system/app)
	ModuleNamespace string      // TS Api 命名空间(SystemManage/AppManage)
	StructRows      []StructRow // 结构体行
	QueryPick       string      // 搜索字段 TS Pick 如 'id' | 'name'
	HasBase         bool
	HasDel          bool   // 是否含有软删除
	HasStatusEum    bool   // 是否有Status枚举(用于生成statusOptions)
	TimeStr         string // 生成时间
}

type StructRow struct {
	Typ       string    // golang类型
	ToTyp     string    // 映射的类型， 如int64映射成string
	TsTyp     string    // TypeScript类型，如int64->string
	Name      string    // 大驼峰
	Query     bool      // 搜索项
	LikeQuery bool      // 模糊搜索(string类型Query字段使用LIKE)
	Must      bool      // 必填项
	Eum       []EumBody // 枚举内容
	Max       string    // 最大长度
	XormTag   string
	JsonTag   string
	QueryTag  string
	Doc       string // 描述(含枚举)
	NameDoc   string // 字段名
}

type EumBody struct {
	Tab string
	Val string
}

type XormTag struct {
	Pk      bool      // 是否主键
	Notnull bool      // 是否必填
	Query   bool      // 是否搜索项
	Varchar string    // 长度
	Default string    // 默认值
	Comment string    // 描述
	Eum     string    // 枚举
	EumBody []EumBody // 枚举体
}

// 构建中间参数，用以存储 xorm、json 内容
type Field struct {
	Typ  string
	Name string
	Xorm string
	Tag  XormTag
	Doc  string
}

type TagOption struct {
	Tags []string
}

func NewTag() *TagOption {
	return &TagOption{
		Tags: make([]string, 0),
	}
}

func (m *TagOption) Len() int {
	return len(m.Tags)
}
func (m *TagOption) WithTag(tag string) *TagOption {
	m.Tags = append(m.Tags, tag)
	return m
}

func (m *TagOption) Builder(sp string) string {
	return strings.Join(m.Tags, sp)
}
