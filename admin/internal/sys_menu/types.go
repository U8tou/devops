package sysmenu

/**
Notes: 系统菜单/权限 API
Time:  2025-01-29
*/

type ListReq struct {
	Current int    `json:"current" default:"1" validate:"required,min=1,max=999999"` // 页码
	Size    int    `json:"size" default:"20" validate:"required,min=1,max=999999"`   // 页长
	Id      string `json:"id" validate:"omitempty"`                                  // ID
	Pid     string `json:"pid" validate:"omitempty"`                                 // 父ID
	Types   int8   `json:"types" validate:"omitempty,oneof=1 2"`                     // 类型:1_菜单,2_权限
	Permis  string `json:"permis" validate:"omitempty,max=50"`                       // 权限代码
}
type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}
type ListBody struct {
	Id         string `json:"id"`         // ID
	Pid        string `json:"pid"`        // 父ID
	Types      int8   `json:"types"`      // 类型:1_菜单,2_权限
	Permis     string `json:"permis"`     // 权限代码
	Remark     string `json:"remark"`     // 备注
	CreateTime string `json:"createTime"` // 创建时间
}

type GetReq struct {
	Id string `json:"id" validate:"required"`
}
type GetResp struct {
	Id         string `json:"id"`         // ID
	Pid        string `json:"pid"`        // 父ID
	Types      int8   `json:"types"`      // 类型:1_菜单,2_权限
	Permis     string `json:"permis"`     // 权限代码
	Remark     string `json:"remark"`     // 备注
	CreateTime string `json:"createTime"` // 创建时间
	CreateBy   string `json:"createBy"`   // 创建者
	UpdateTime string `json:"updateTime"` // 更新时间
	UpdateBy   string `json:"updateBy"`   // 更新者
}

type DelReq struct {
	Ids []string `json:"ids" validate:"required,dive,number"` // IDs
}
type DelResp struct {
	Affect string `json:"affect"` // 影响条数
}

type AddReq struct {
	Pid    string `json:"pid"`                                             // 父ID
	Types  int8   `json:"types" default:"1" validate:"required,oneof=1 2"` // 类型:1_菜单,2_权限
	Permis string `json:"permis" validate:"required,max=50"`               // 权限代码
	Remark string `json:"remark" validate:"omitempty,max=50"`              // 备注
}
type AddResp struct {
	Id string `json:"id"` // 新增ID
}

type EditReq struct {
	Id     string `json:"id" validate:"required"`                          // ID
	Pid    string `json:"pid"`                                             // 父ID
	Types  int8   `json:"types" default:"1" validate:"required,oneof=1 2"` // 类型:1_菜单,2_权限
	Permis string `json:"permis" validate:"required,max=50"`               // 权限代码
	Remark string `json:"remark" validate:"omitempty,max=50"`              // 备注
}
type EditResp struct {
	Affect string `json:"affect"` // 影响条数
}
