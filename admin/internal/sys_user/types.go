package sysuser

/**
Notes: 系统用户 API
Time:  2025-12-18 17:46:03
*/

type ListReq struct {
	Current   int      `json:"current" default:"1" validate:"required,min=1,max=999999"` // 页码
	Size      int      `json:"size" default:"20" validate:"required,min=1,max=999999"`   // 页长
	Id        string   `json:"id" validate:"omitempty,numeric"`                          // 用户ID
	UserName  string   `json:"userName" validate:"omitempty,max=30"`                     // 用户账号
	NickName  string   `json:"nickName" validate:"omitempty,max=30"`                     // 用户昵称
	Email     string   `json:"email" validate:"omitempty,email,max=50"`                  // 用户邮箱
	Phone     string   `json:"phone" validate:"omitempty,max=20"`                        // 电话号码
	Sex       int8     `json:"sex" default:"1" validate:"omitempty,oneof=1 2 3"`         // 用户性别:1_男,2_女,3_未知
	Status    int8     `json:"status" default:"1" validate:"omitempty,max=1,oneof=1 2"`  // 状态:1_正常,2_停用
	DeptId    string   `json:"deptId" validate:"omitempty,numeric"`                      // 部门ID
	TimeRange []string `json:"timeRange" validate:"omitempty,len=2,dive,numeric"`        // 时间范围（创建时间，Unix 秒时间戳区间：[开始,结束]）
}
type ListResp struct {
	Total string     `json:"total"`
	Rows  []ListBody `json:"rows"`
}
type ListBody struct {
	Id         string   `json:"id" validate:"required,numeric"`                           // 用户ID
	UserName   string   `json:"userName" validate:"required,max=30"`                      // 用户账号
	NickName   string   `json:"nickName" validate:"required,max=30"`                      // 用户昵称
	UserType   string   `json:"userType" default:"00" validate:"required,max=2,oneof=00"` // 用户类型:00_系统用户
	Email      string   `json:"email" validate:"required,email,max=50"`                   // 用户邮箱
	PhoneArea  string   `json:"phoneArea" default:"+86" validate:"required,max=10"`       // 电话区号
	Phone      string   `json:"phone" validate:"required,max=20"`                         // 电话号码
	Sex        int8     `json:"sex" default:"1" validate:"required,oneof=1 2 3"`          // 用户性别:1_男,2_女,3_未知
	Avatar     string   `json:"avatar" validate:"required,max=100"`                       // 头像
	Status     int8     `json:"status" default:"1" validate:"required,max=1,oneof=1 2"`   // 状态:1_正常,2_停用
	Address    string   `json:"address" validate:"required,max=255"`                      // 联系地址
	Remark     string   `json:"remark" validate:"required,max=500"`                       // 备注
	CreateTime string   `json:"createTime" validate:"required"`                           // 创建时间
	CreateBy   string   `json:"createBy" validate:"required"`                             // 创建者
	UpdateTime string   `json:"updateTime" validate:"required"`                           // 更新时间
	UpdateBy   string   `json:"updateBy" validate:"required"`                             // 更新者
	Depts      []string `json:"depts"`                                                    // 关联部门IDs
	Roles      []string `json:"roles"`                                                    // 关联角色IDs
	Posts      []string `json:"posts"`                                                    // 关联岗位IDs
}

type GetReq struct {
	Id string `json:"id" validate:"required,number"`
}
type GetResp struct {
	Id         string   `json:"id" validate:"required,numeric"`                           // 用户ID
	UserName   string   `json:"userName" validate:"required,max=30"`                      // 用户账号
	NickName   string   `json:"nickName" validate:"required,max=30"`                      // 用户昵称
	UserType   string   `json:"userType" default:"00" validate:"required,max=2,oneof=00"` // 用户类型:00_系统用户
	Email      string   `json:"email" validate:"required,email,max=50"`                   // 用户邮箱
	PhoneArea  string   `json:"phoneArea" default:"+86" validate:"required,max=10"`       // 电话区号
	Phone      string   `json:"phone" validate:"required,max=20"`                         // 电话号码
	Sex        int8     `json:"sex" default:"1" validate:"required,oneof=1 2 3"`          // 用户性别:1_男,2_女,3_未知
	Avatar     string   `json:"avatar" validate:"required,max=100"`                       // 头像
	Status     int8     `json:"status" default:"1" validate:"required,max=1,oneof=1 2"`   // 状态:1_正常,2_停用
	Address    string   `json:"address" validate:"required,max=255"`                      // 联系地址
	Remark     string   `json:"remark" validate:"required,max=500"`                       // 备注
	CreateTime string   `json:"createTime" validate:"required"`                           // 创建时间
	CreateBy   string   `json:"createBy" validate:"required"`                             // 创建者
	UpdateTime string   `json:"updateTime" validate:"required"`                           // 更新时间
	UpdateBy   string   `json:"updateBy" validate:"required"`                             // 更新者
	Depts      []string `json:"depts"`                                                    // 关联部门IDs
	Roles      []string `json:"roles"`                                                    // 关联角色IDs
	Posts      []string `json:"posts"`                                                    // 关联岗位IDs
}

type DelReq struct {
	Ids []string `json:"ids" validate:"required,dive,number"` // IDs
}
type DelResp struct {
	Affect string `json:"affect" validate:"required"` // 影响条数
}

type AddReq struct {
	UserName  string `json:"userName" validate:"required,max=30"`                    // 用户账号
	NickName  string `json:"nickName" validate:"required,max=30"`                    // 用户昵称
	Email     string `json:"email" validate:"omitempty,email,max=50"`                 // 用户邮箱（选填）
	PhoneArea string `json:"phoneArea" default:"+86" validate:"required,max=10"`     // 电话区号
	Phone     string `json:"phone" validate:"required,max=20"`                       // 电话号码
	Sex       int8   `json:"sex" default:"1" validate:"required,oneof=1 2 3"`        // 用户性别:1_男,2_女,3_未知
	Avatar    string `json:"avatar" validate:"omitempty,max=100"`                    // 头像
	Password  string `json:"password" validate:"required,min=6,max=100"`             // 密码
	Status    int8   `json:"status" default:"1" validate:"required,max=1,oneof=1 2"` // 状态:1_正常,2_停用
	Address   string `json:"address" validate:"omitempty,max=255"`                   // 联系地址
	Remark    string `json:"remark" validate:"omitempty,max=500"`                    // 备注
}
type AddResp struct {
	Affect string `json:"affect"` // 影响条数
}

type EditReq struct {
	Id        string `json:"id" validate:"required,numeric"`                         // 用户ID
	UserName  string `json:"userName" validate:"required,max=30"`                    // 用户账号
	NickName  string `json:"nickName" validate:"required,max=30"`                    // 用户昵称
	Email     string `json:"email" validate:"omitempty,email,max=50"`                 // 用户邮箱（选填）
	PhoneArea string `json:"phoneArea" default:"+86" validate:"required,max=10"`     // 电话区号
	Phone     string `json:"phone" validate:"required,max=20"`                       // 电话号码
	Sex       int8   `json:"sex" default:"1" validate:"required,oneof=1 2 3"`        // 用户性别:1_男,2_女,3_未知
	Avatar    string `json:"avatar" validate:"omitempty,max=100"`                    // 头像
	Password  string `json:"password" validate:"omitempty,min=6,max=100"`            // 密码，空表示不修改
	Status    int8   `json:"status" default:"1" validate:"required,max=1,oneof=1 2"` // 状态:1_正常,2_停用
	Address   string `json:"address" validate:"omitempty,max=255"`                   // 联系地址
	Remark    string `json:"remark" validate:"omitempty,max=500"`                    // 备注
}
type EditResp struct {
	Affect string `json:"affect"` // 影响条数
}

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	UserId   string `json:"userId" validate:"required,number"`          // 用户ID
	Password string `json:"password" validate:"required,min=6,max=100"` // 新密码
}

// ResetPasswordResp 重置密码响应
type ResetPasswordResp struct {
	Affect string `json:"affect"` // 影响条数
}

// AssignRoleReq 分配角色请求
type AssignRoleReq struct {
	UserId  string   `json:"userId" validate:"required,number"`       // 用户ID
	RoleIds []string `json:"roleIds" validate:"required,dive,number"` // 角色ID列表
}

// AssignRoleResp 分配角色响应
type AssignRoleResp struct {
	Affect string `json:"affect"` // 影响条数
}

// AssignDeptReq 分配部门请求
type AssignDeptReq struct {
	UserId  string   `json:"userId" validate:"required,number"`       // 用户ID
	DeptIds []string `json:"deptIds" validate:"required,dive,number"` // 部门ID列表
}

// AssignDeptResp 分配部门响应
type AssignDeptResp struct {
	Affect string `json:"affect"` // 影响条数
}

// AssignPostReq 分配岗位请求
type AssignPostReq struct {
	UserId  string   `json:"userId" validate:"required,number"`       // 用户ID
	PostIds []string `json:"postIds" validate:"required,dive,number"` // 岗位ID列表
}

// AssignPostResp 分配岗位响应
type AssignPostResp struct {
	Affect string `json:"affect"` // 影响条数
}

// ImportResp 用户导入响应
type ImportResp struct {
	SuccessCount int          `json:"successCount"` // 成功条数
	FailList     []ImportFail `json:"failList"`     // 失败行（行号+原因）
}

// ImportFail 导入失败行
type ImportFail struct {
	Row    int    `json:"row"`    // Excel 行号（从 2 起为数据行）
	Reason string `json:"reason"` // 失败原因
}
