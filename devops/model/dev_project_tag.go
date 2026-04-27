package model

// DevProjectTag 项目标签字典（id 与项目表一致使用雪花）
type DevProjectTag struct {
	Id   int64  `xorm:"pk notnull comment('ID')"`
	Name string `xorm:"varchar(128) notnull unique(dev_project_tag_name) comment('标签名')"`
}

func (m *DevProjectTag) Comment() string {
	return "项目标签"
}

// DevProjectTagLink 项目与标签多对多；删除字典标签时不级联删本表，以支持「其他」孤儿关联
type DevProjectTagLink struct {
	ProjectId int64 `xorm:"pk notnull comment('项目ID')"`
	TagId     int64 `xorm:"pk notnull comment('标签ID')"`
}

func (m *DevProjectTagLink) Comment() string {
	return "项目标签关联"
}

// DevProjectTagBrief 列表/详情展示用（孤儿标签 name 为空）
type DevProjectTagBrief struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
