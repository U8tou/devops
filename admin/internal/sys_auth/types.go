package sysauth

type LoginReq struct {
	UserName string `json:"userName" example:"admin" validate:"required,min=2,max=20"`  // 用户名
	Password string `json:"password" example:"123456" validate:"required,min=5,max=20"` // 密码
	CodeId   string `json:"codeId"`
	Code     string `json:"code"`
}

type LoginResp struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// RefreshTokenReq Token刷新请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refreshToken" example:"" validate:"required"` // 刷新令牌
}
type InfoResp struct {
	UserId     string   `json:"userId" validate:"required,numeric"`                       // 用户ID
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
	Buttons    []string `json:"buttons"`                                                  // 按钮权限
	Menus      []string `json:"menus"`                                                    // 菜单权限
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" example:"123456" validate:"required,min=5,max=20"` // 原密码
	NewPassword string `json:"newPassword" example:"654321" validate:"required,min=5,max=20"` // 新密码
}

// UpdateProfileReq 个人信息更新请求（仅允许修改的字段）
type UpdateProfileReq struct {
	NickName  string `json:"nickName" validate:"omitempty,max=30"`    // 用户昵称
	Email     string `json:"email" validate:"omitempty,email,max=50"` // 用户邮箱
	PhoneArea string `json:"phoneArea" validate:"omitempty,max=10"`   // 电话区号
	Phone     string `json:"phone" validate:"omitempty,max=20"`       // 电话号码
	Sex       int8   `json:"sex" validate:"omitempty,oneof=1 2 3"`    // 用户性别:1_男,2_女,3_未知
	Address   string `json:"address" validate:"omitempty,max=255"`    // 联系地址
}

// UploadAvatarResp 头像上传响应
type UploadAvatarResp struct {
	Avatar string `json:"avatar"` // 头像地址（相对路径或完整 URL）
}

// RegisterReq 用户注册请求
type RegisterReq struct {
	UserName string `json:"userName" validate:"required,min=3,max=20"` // 用户名
	Password string `json:"password" validate:"required,min=6,max=20"` // 密码
}

// RegisterResp 用户注册响应
type RegisterResp struct {
	UserName string `json:"userName"` // 注册成功的用户名
}
