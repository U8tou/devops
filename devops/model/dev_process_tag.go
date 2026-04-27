package model

// DevProcessTag 流程标签字典（仅 id、name）；id 与流程表一致使用雪花
type DevProcessTag struct {
	Id   int64  `xorm:"pk notnull comment('ID')"`
	Name string `xorm:"varchar(128) notnull unique(dev_process_tag_name) comment('标签名')"`
}

func (m *DevProcessTag) Comment() string {
	return "流程标签"
}

// DevProcessTagLink 流程与标签多对多；删除字典标签时不级联删本表，以支持「其他」孤儿关联
type DevProcessTagLink struct {
	ProcessId int64 `xorm:"pk notnull comment('流程ID')"`
	TagId     int64 `xorm:"pk notnull comment('标签ID')"`
}

func (m *DevProcessTagLink) Comment() string {
	return "流程标签关联"
}

// DevProcessTagBrief 列表/详情展示用（孤儿标签 name 为空）
type DevProcessTagBrief struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
