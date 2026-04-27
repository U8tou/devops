package syspost

type ListReq struct {
	Current int    `json:"current" default:"1" validate:"required,min=1,max=999999"`
	Size    int    `json:"size" default:"20" validate:"required,min=1,max=999999"`
	Id      string `json:"id"`
	Name    string `json:"name"`
	Sort    int    `json:"sort"`
	Status  int8   `json:"status,omitempty" validate:"omitempty,oneof=1 2"`
}

type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}

type ListBody struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Sort       int    `json:"sort"`
	Status     int8   `json:"status"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type GetReq struct {
	Id string `json:"id" validate:"required,number"`
}

type GetResp struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Sort       int    `json:"sort"`
	Status     int8   `json:"status"`
	CreateTime string `json:"createTime"`
	CreateBy   string `json:"createBy"`
	UpdateTime string `json:"updateTime"`
	UpdateBy   string `json:"updateBy"`
}

type DelReq struct {
	Ids []string `json:"ids" validate:"required,dive,number"`
}

type DelResp struct {
	Affect string `json:"affect"`
}

type AddReq struct {
	Name   string `json:"name" validate:"required,max=50"`
	Sort   int    `json:"sort" validate:"required"`
	Status int8   `json:"status,omitempty" validate:"omitempty,oneof=1 2"`
}

type AddResp struct {
	Affect string `json:"affect"`
}

type EditReq struct {
	Id     string `json:"id" validate:"required,number"`
	Name   string `json:"name" validate:"required,max=50"`
	Sort   int    `json:"sort" validate:"required"`
	Status int8   `json:"status,omitempty" validate:"omitempty,oneof=1 2"`
}

type EditResp struct {
	Affect string `json:"affect"`
}
