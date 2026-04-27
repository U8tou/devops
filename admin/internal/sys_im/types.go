package sysim

// ImClient IM 聊天室客户端
type ImClient struct {
	ID        string
	NickName  string
	WriteChan chan []byte
}

// clientMsg 客户端 -> 服务端
type clientMsg struct {
	Type    string `json:"type"`    // join, leave, message
	RoomId  string `json:"roomId"`
	Content string `json:"content"`
}

// serverMsg 服务端 -> 客户端
// type 枚举及字段说明：
// - message: 聊天消息。含 roomId,from,loginId,connId,nickName,avatar,content,time
// - joined:  加入成功确认。含 roomId,memberCount
// - left:    离开成功确认。含 roomId,memberCount
// - error:   错误提示。含 msg
// - online:  成员上线广播。含 roomId,from,loginId,connId,nickName,avatar,msg,time,memberCount
// - offline: 成员下线广播。含 roomId,from,loginId,connId,nickName,avatar,msg,time,memberCount
type serverMsg struct {
	Type     string `json:"type"`               // message|joined|left|error|online|offline
	RoomId   string `json:"roomId,omitempty"`   // 房间ID
	From     string `json:"from,omitempty"`     // 会话标识 loginId_connId
	LoginId  string `json:"loginId,omitempty"`  // 登录ID，游客为 0
	ConnId   string `json:"connId,omitempty"`   // 连接ID，分布式 ID 生成的唯一标识
	NickName string `json:"nickName,omitempty"` // 用户昵称，如 小王、游客
	Avatar   string `json:"avatar,omitempty"`   // 头像完整 URL
	Content  string `json:"content,omitempty"`  // 消息内容
	Time         int64  `json:"time,omitempty"`         // 时间戳毫秒
	Msg          string `json:"msg,omitempty"`           // 提示文案，如 "xxx 上线了"、"xxx 下线了"、错误信息
	MemberCount  int    `json:"memberCount,omitempty"`  // 房间当前在线人数，joined/left/online/offline 时推送
}
