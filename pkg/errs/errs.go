package errs

import (
	"fmt"
)

// 自定义错误类型（可选）
type Errs struct {
	Code int    // 错误代码
	Msg  string // 错误消息
	Data any    // 错误内容
}

func (e *Errs) Error() string {
	return fmt.Sprintf("%d:%s", e.Code, e.Msg)
}

func errMsg(code int, msg string) *Errs {
	return &Errs{
		Code: code,
		Msg:  msg,
	}
}

// 自定义异常
func New(msg string) *Errs {
	return &Errs{
		Code: ERR_SYS_CUSTOM.Code,
		Msg:  msg,
	}
}

// Args 参数异常
func Args(err error) *Errs {
	return &Errs{
		Code: ERR_SYS_ARGS.Code,
		Msg:  ERR_SYS_ARGS.Msg,
		Data: err.Error(),
	}
}

// Sys 系统内部异常
// 注意此處返回的是error，而不是 *Errs
// 如果返回*Errs 會導致 err != nil 的判斷失效
func Sys(err error) error {
	if err == nil {
		return nil
	}
	return &Errs{
		Code: ERR_SYS_SERVER.Code,
		Msg:  ERR_SYS_SERVER.Msg,
		Data: err.Error(),
	}
}

// 异常
var (
	CODE_400 = 400 // 请求异常
	CODE_401 = 401 // 认证失败
	CODE_403 = 403 // 权限异常
	CODE_404 = 404 // 不存在
	CODE_405 = 405 // 处理异常
	CODE_500 = 500 // 服务异常

	// 系统错误
	ERR_SYS_ARGS = errMsg(400, "参数错误")
	ERR_SYS_401  = errMsg(401, "认证失败,无法访问系统资源")
	ERR_SYS_403  = errMsg(403, "当前操作没有权限")
	ERR_SYS_404  = errMsg(404, "访问资源不存在")
	ERR_SYS_405  = errMsg(405, "操作失败,请检查请求内容是否正确") // 系统应友好提示
	ERR_SYS_500  = errMsg(500, "系统错误，请反馈给管理员")

	ERR_SYS_DEFAULT = errMsg(500000, "未知异常")
	ERR_SYS_SERVER  = errMsg(500100, "服务异常")
	ERR_SYS_CUSTOM  = errMsg(500101, "自定义异常")
	ERR_SYS_OPER    = errMsg(500102, "操作失败")

	ERR_LOGIN        = errMsg(500103, "账号或密码错误")
	ERR_LOGINOUT     = errMsg(500104, "登出失败")
	ERR_LOGIN_VERIFY = errMsg(500105, "验证码验证失败")
	ERR_LOGIN_LOCK   = errMsg(500106, "账号已锁定")
	ERR_LOGIN_AUTH   = errMsg(500107, "Token认证失败")
	ERR_NOT_ACCOUNT  = errMsg(500108, "账号不存在")
	ERR_PWD          = errMsg(500109, "密码错误")
	ERR_LOGIN_ARGS   = errMsg(500110, "账号或密码错误")
	ERR_HAS_EXIST    = errMsg(500111, "已存在")
	ERR_NOT_EXIST    = errMsg(500112, "不存在")
	ERR_NOT_SAME     = errMsg(500113, "新旧密码不能一致")
	ERR_HAS_ACCOUNT  = errMsg(500114, "已存在账号")
	ERR_HAS_PHONE    = errMsg(500115, "已存在号码")
	ERR_HAS_EMAIL    = errMsg(500116, "已存在邮箱")
	ERR_ACCOUNT_CHAR = errMsg(500117, "账号包含非法字符")
	ERR_NO_ROLE      = errMsg(500118, "账号尚未分配角色，请联系管理员")
	// 数据庫异常
	ERR_DB_OPEN     = errMsg(500201, "数据库连接失败")
	ERR_DB_READ     = errMsg(500202, "数据获取失败")
	ERR_DB_EXIST    = errMsg(500203, "数据已存在")
	ERR_DB_NO_EXIST = errMsg(500204, "数据不存在")
	// 文件错误
	ERR_FILE_OPEN  = errMsg(500301, "文件打开失败")
	ERR_FILE_REDE  = errMsg(500302, "文件读取失败")
	ERR_FILE_WRITE = errMsg(500303, "文件写入失败")

	// 业务错误
	ERR_SEN_ONUSE = errMsg(500501, "该时间段正在使用或已被预订")
	ERR_SEN_STORE = errMsg(500502, "库存操作失败")
	ERR_SEN_ROOM  = errMsg(500503, "房间不可用")
)
