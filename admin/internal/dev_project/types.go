package dev_project

type ListReq struct {
	Current int    `query:"current" default:"1" validate:"required,min=1,max=999999"`
	Size    int    `query:"size" default:"20" validate:"required,min=1,max=999999"`
	Name    string `query:"name"`
	Status  string `query:"status"`
	TagIds  string `query:"tagIds"`
	TagOther string `query:"tagOther"`
	TagMode string `query:"tagMode"`
}

type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}

type ListBody struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Progress         string    `json:"progress"`
	VersionChangelog string    `json:"versionChangelog"`
	CreateTime       string    `json:"createTime"`
	UpdateTime       string    `json:"updateTime"`
	Tags             []TagItem `json:"tags"`
}

// TagItem 项目标签（name 为空表示字典已删除，仅孤儿关联）
type TagItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetReq struct {
	Id string `query:"id" validate:"required,number"`
}

type GetResp struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Progress         string    `json:"progress"`
	VersionChangelog string    `json:"versionChangelog"`
	MindJson         string    `json:"mindJson"`
	CreateTime       string    `json:"createTime"`
	CreateBy         string    `json:"createBy"`
	UpdateTime       string    `json:"updateTime"`
	UpdateBy         string    `json:"updateBy"`
	Tags             []TagItem `json:"tags"`
}

type AddReq struct {
	Name             string  `json:"name" validate:"required,max=200"`
	Status           int8    `json:"status"`
	Progress         int8    `json:"progress"`
	VersionChangelog string  `json:"versionChangelog"`
	MindJson         string  `json:"mindJson"`
	TagIds           []int64 `json:"tagIds"`
}

type AddResp struct {
	Affect string `json:"affect"`
	Id     string `json:"id"`
}

type EditReq struct {
	Id               string  `json:"id" validate:"required,number"`
	Name             string  `json:"name" validate:"required,max=200"`
	Status           int8    `json:"status"`
	Progress         int8    `json:"progress"`
	VersionChangelog string  `json:"versionChangelog"`
	TagIds           []int64 `json:"tagIds"`
}

type EditResp struct {
	Affect string `json:"affect"`
}

type EditMindReq struct {
	Id       string `json:"id" validate:"required,number"`
	MindJson string `json:"mindJson" validate:"required"`
}

type EditMindResp struct {
	Affect string `json:"affect"`
}

// —— 标签字典（dev_project_tag）——

type TagListResp struct {
	Rows []TagRow `json:"rows"`
}

type TagRow struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type TagAddReq struct {
	Name string `json:"name" validate:"required"`
}

type TagAddResp struct {
	Id string `json:"id"`
}

type TagEditReq struct {
	Id   string `json:"id" validate:"required,number"`
	Name string `json:"name" validate:"required"`
}

type TagDelReq struct {
	Id string `json:"id" query:"id" validate:"required,number"`
}
