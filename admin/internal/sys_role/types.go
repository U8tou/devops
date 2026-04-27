package sysrole

/**
Notes: 系统角色 API
Time:  2025-01-29
*/

type ListReq struct {
	Current int    `json:"current" default:"1" validate:"required,min=1,max=999999"` // 页码
	Size    int    `json:"size" default:"20" validate:"required,min=1,max=999999"`   // 页长
	Id      string `json:"id" validate:"omitempty"`                                  // 角色ID
	Name    string `json:"name" validate:"omitempty,max=30"`                         // 角色名称
	Role    string `json:"role" validate:"omitempty,max=30"`                         // 角色标识
	Status  int8   `json:"status" validate:"omitempty,oneof=1 2"`                    // 状态:1_正常,2_停用
}
type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}
type ListBody struct {
	Id          string `json:"id"`          // 角色ID
	Name        string `json:"name"`        // 名称
	Role        string `json:"role"`        // 角色标识
	Status      int8   `json:"status"`      // 状态:1_正常,2_停用
	MenuLinkage int8   `json:"menuLinkage"` // 操作父子联动:1_联动,2_不联动
	DeptLinkage int8   `json:"deptLinkage"` // 数据父子联动:1_联动,2_不联动
	Sort        int    `json:"sort"`        // 显示顺序
	Remark      string `json:"remark"`      // 备注
	CreateTime  string `json:"createTime"`  // 创建时间
	UpdateTime  string `json:"updateTime"`  // 更新时间
}

type GetReq struct {
	Id string `json:"id" validate:"required"`
}
type GetResp struct {
	Id          string   `json:"id"`          // 角色ID
	Name        string   `json:"name"`        // 名称
	Role        string   `json:"role"`        // 角色标识
	Status      int8     `json:"status"`      // 状态:1_正常,2_停用
	MenuLinkage int8     `json:"menuLinkage"` // 操作父子联动
	DeptLinkage int8     `json:"deptLinkage"` // 数据父子联动
	Sort        int      `json:"sort"`        // 显示顺序
	Remark      string   `json:"remark"`      // 备注
	MenuIds     []string `json:"menuIds"`     // 菜单ID列表
	DeptIds     []string `json:"deptIds"`     // 部门ID列表
	CreateTime  string   `json:"createTime"`  // 创建时间
	CreateBy    string   `json:"createBy"`    // 创建者
	UpdateTime  string   `json:"updateTime"`  // 更新时间
	UpdateBy    string   `json:"updateBy"`    // 更新者
}

type DelReq struct {
	Ids []string `json:"ids" validate:"required,dive"` // IDs
}
type DelResp struct {
	Affect string `json:"affect"` // 影响条数
}

type AddReq struct {
	Name   string `json:"name" validate:"required,max=30"`         // 名称
	Role   string `json:"role" validate:"required,max=30"`         // 角色标识
	Status int8   `json:"status" default:"1" validate:"oneof=1 2"` // 状态:1_正常,2_停用
	Sort   int    `json:"sort" default:"999"`                      // 显示顺序
	Remark string `json:"remark" validate:"omitempty,max=500"`     // 备注
}
type AddResp struct {
	Id string `json:"id"` // 新增ID
}

type EditReq struct {
	Id     string `json:"id" validate:"required"`                  // 角色ID
	Name   string `json:"name" validate:"required,max=30"`         // 名称
	Role   string `json:"role" validate:"required,max=30"`         // 角色标识
	Status int8   `json:"status" default:"1" validate:"oneof=1 2"` // 状态:1_正常,2_停用
	Sort   int    `json:"sort" default:"999"`                      // 显示顺序
	Remark string `json:"remark" validate:"omitempty,max=500"`     // 备注
}
type EditResp struct {
	Affect string `json:"affect"` // 影响条数
}

type AssignMenuReq struct {
	RoleId      string   `json:"roleId" validate:"required,number"`       // 角色ID
	MenuLinkage int8     `json:"menuLinkage" validate:"required,number"`  // 操作父子联动:1_联动,2_不联动
	MenuIds     []string `json:"menuIds" validate:"required,dive,number"` // 菜单ID列表
}
type AssignMenuResp struct {
	Affect string `json:"affect"` // 影响条数
}

type AssignDeptReq struct {
	RoleId      string   `json:"roleId" validate:"required,number"`       // 角色ID
	DeptLinkage int8     `json:"deptLinkage" validate:"required,number"`  // 数据父子联动:1_联动,2_不联动
	DeptIds     []string `json:"deptIds" validate:"required,dive,number"` // 部门ID列表
}
type AssignDeptResp struct {
	Affect string `json:"affect"` // 影响条数
}
