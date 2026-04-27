package sysdept

type ListReq struct {
	Current int    `json:"current" default:"1" validate:"required,min=1,max=999999"` // 页码
	Size    int    `json:"size" default:"20" validate:"required,min=1,max=999999"`   // 页长
	Id      string `json:"id" `                                                      // ID
	Pid     string `json:"pid" `                                                     // 上级部门ID
	Name    string `json:"name" `                                                    // 部门名称
	Profile string `json:"profile,omitempty" `                                       // 部门简介
	Leader  string `json:"leader" `                                                  // 负责人
	Phone   string `json:"phone" `                                                   // 负责人电话
	Email   string `json:"email" `                                                   // 负责人邮箱
	Sort    int    `json:"sort" `                                                    // 排序
	Status  int8   `json:"status,omitempty" validate:"omitempty,oneof=1 2"`          // 状态:1_正常,0_停用
}

type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}

type ListBody struct {
	Id         string `json:"id" validate:"required"`                // ID
	Pid        string `json:"pid" `                                  // 上级部门ID
	Name       string `json:"name" validate:"required"`              // 部门名称
	Profile    string `json:"profile,omitempty" `                    // 部门简介
	Leader     string `json:"leader" validate:"required"`            // 负责人
	Phone      string `json:"phone" validate:"required"`             // 负责人电话
	Email      string `json:"email" validate:"required"`             // 负责人邮箱
	Sort       int    `json:"sort" validate:"required"`              // 排序
	Status     int8   `json:"status,omitempty" validate:"oneof=1 2"` // 状态:1_正常,2_停用
	CreateTime string `json:"createTime"`                            // 创建时间
	UpdateTime string `json:"updateTime"`                            // 更新时间
}

type GetReq struct {
	Id string `json:"id" validate:"required,number"`
}

type GetResp struct {
	Id         string `json:"id" validate:"required,number"`         // ID
	Pid        string `json:"pid" `                                  // 上级部门ID
	Name       string `json:"name" validate:"required"`              // 部门名称
	Profile    string `json:"profile,omitempty" `                    // 部门简介
	Leader     string `json:"leader" validate:"required"`            // 负责人
	Phone      string `json:"phone" validate:"required"`             // 负责人电话
	Email      string `json:"email" validate:"required"`             // 负责人邮箱
	Sort       int    `json:"sort" validate:"required"`              // 排序
	Status     int8   `json:"status,omitempty" validate:"oneof=1 2"` // 状态:1_正常,2_停用
	CreateTime string `json:"createTime"`                            // 创建时间
	CreateBy   string `json:"createBy"`                              // 创建者
	UpdateTime string `json:"updateTime"`                            // 更新时间
	UpdateBy   string `json:"updateBy"`                              // 更新者
}

type DelReq struct {
	Ids []string `json:"ids" validate:"required,dive,number"` // IDs
}

type DelResp struct {
	Affect string `json:"affect" validate:"required"` // 影响条数
}

type AddReq struct {
	Pid     string `json:"pid" `                                  // 上级部门ID
	Name    string `json:"name" validate:"required"`              // 部门名称
	Profile string `json:"profile,omitempty" `                    // 部门简介
	Leader  string `json:"leader" validate:"required"`            // 负责人
	Phone   string `json:"phone"`                                 // 电话
	Email   string `json:"email"`                                 // 邮箱
	Sort    int    `json:"sort" validate:"required"`              // 排序
	Status  int8   `json:"status,omitempty" validate:"oneof=1 2"` // 状态:1_正常,2_停用
}

type AddResp struct {
	Affect string `json:"affect"` // 影响条数
}

type EditReq struct {
	Id      string `json:"id" validate:"required"`                // ID
	Pid     string `json:"pid" `                                  // 上级部门ID
	Name    string `json:"name" validate:"required"`              // 部门名称
	Profile string `json:"profile,omitempty" `                    // 部门简介
	Leader  string `json:"leader" validate:"required"`            // 负责人
	Phone   string `json:"phone"`                                 // 电话
	Email   string `json:"email"`                                 // 邮箱
	Sort    int    `json:"sort" validate:"required"`              // 排序
	Status  int8   `json:"status,omitempty" validate:"oneof=1 2"` // 状态:1_正常,0_停用
}

type EditResp struct {
	Affect string `json:"affect"` // 影响条数
}
