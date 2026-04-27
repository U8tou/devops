package dev_process

type ListReq struct {
	Current int    `json:"current" default:"1" validate:"required,min=1,max=999999"`
	Size    int    `json:"size" default:"20" validate:"required,min=1,max=999999"`
	Id      string `json:"id"`
	Code    string `json:"code"`
	Remark  string `json:"remark"`
	// tagIds 逗号分隔的标签 id；与 tagOther 组合为 OR 筛选
	TagIds string `query:"tagIds"`
	// tagOther 为 1/true 时包含「孤儿标签」关联的流程
	TagOther string `query:"tagOther"`
	// tagMode：include（默认）包含所选条件；exclude 不包含（条件取反）
	TagMode string `query:"tagMode"`
}

type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}

type ListBody struct {
	Id                 string `json:"id"`
	Code               string `json:"code"`
	Remark             string `json:"remark"`
	CronExpr           string `json:"cronExpr"`
	CronEnabled        string `json:"cronEnabled"`
	LastExecTime       string `json:"lastExecTime"`
	LastExecDurationMs string `json:"lastExecDurationMs"`
	LastExecResult     string    `json:"lastExecResult"`
	CreateTime         string    `json:"createTime"`
	UpdateTime         string    `json:"updateTime"`
	Tags               []TagItem `json:"tags"`
}

// TagItem 流程标签（name 为空表示字典已删除，仅孤儿关联）
type TagItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetReq struct {
	Id string `json:"id" validate:"required,number"`
}

type GetResp struct {
	Id                 string            `json:"id"`
	Code               string            `json:"code"`
	Remark             string            `json:"remark"`
	Flow               string            `json:"flow"`
	Env                map[string]string `json:"env"`
	CronExpr           string            `json:"cronExpr"`
	CronEnabled        string            `json:"cronEnabled"`
	LastExecTime       string            `json:"lastExecTime"`
	LastExecDurationMs string            `json:"lastExecDurationMs"`
	LastExecResult     string            `json:"lastExecResult"`
	LastExecLog        string            `json:"lastExecLog"`
	CreateTime         string            `json:"createTime"`
	CreateBy           string            `json:"createBy"`
	UpdateTime         string            `json:"updateTime"`
	UpdateBy           string            `json:"updateBy"`
	Tags               []TagItem         `json:"tags"`
}

type AddReq struct {
	Code        string            `json:"code" validate:"required,max=64"`
	Remark      string            `json:"remark" validate:"max=500"`
	Flow        string            `json:"flow"`
	Env         map[string]string `json:"env"`
	CronEnabled int8              `json:"cronEnabled"`
	CronExpr    string            `json:"cronExpr"`
	TagIds      []int64           `json:"tagIds"`
}

type AddResp struct {
	Affect string `json:"affect"`
	Id     string `json:"id"`
}

type EditReq struct {
	Id          string  `json:"id" validate:"required,number"`
	Code        string  `json:"code" validate:"required,max=64"`
	Remark      string  `json:"remark" validate:"max=500"`
	CronEnabled int8    `json:"cronEnabled"`
	CronExpr    string  `json:"cronExpr"`
	TagIds      []int64 `json:"tagIds"`
}

type EditResp struct {
	Affect string `json:"affect"`
}

type EditFlowReq struct {
	Id   string `json:"id" validate:"required,number"`
	Flow string `json:"flow"`
}

type EditFlowResp struct {
	Affect string `json:"affect"`
}

// —— 标签字典（dev_process_tag）——

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
