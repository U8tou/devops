package sysim

import (
	"context"
	"encoding/json"
	"log/slog"
	"pkg/auth"
	"pkg/conf"
	"strconv"
	"strings"
	"sync"
	"time"

	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/yitter/idgenerator-go/idgen"
)

const (
	writeChanSize = 100
)

// Upgrader 返回 Fiber WebSocket 配置
func Upgrader() fiberws.Config {
	return fiberws.Config{
		RecoverHandler: func(conn *fiberws.Conn) {
			if err := recover(); err != nil {
				slog.Error("sys_im WebSocket panic", "err", err)
			}
		},
	}
}

// Handler IM 聊天室 WebSocket 处理器
//
// @Summary		IM 聊天室 WebSocket
// @Description	建立 IM 聊天室 WebSocket 连接。可选传 token 鉴权：有 token 且有效则显示登录用户昵称和头像，无 token 为游客。nickName 为纯昵称如 小王、游客；会话唯一性由 loginId、connId 区分。连接后通过 JSON 消息交互：<br/>客户端→服务端：join(roomId)、leave(roomId)、message(roomId,content)；<br/>服务端→客户端：message、joined、left、error、online(成员上线广播)、offline(成员下线广播)
// @Tags			SysImApi
// @Accept			json
// @Produce		json
// @Param			token	query		string	false	"鉴权 token（可选），有则显示用户昵称，无则为游客"
// @Success		101		{string}	string	"Switching Protocols，升级为 WebSocket"
// @Failure		400		{string}	string	"非 WebSocket 升级请求"
// @Router			/sys_im/chat [get]
func Handler(conn *fiberws.Conn) {
	ctx := conn.Locals("ctx")
	if ctx == nil {
		ctx = context.Background()
	}
	cctx, ok := ctx.(context.Context)
	if !ok {
		cctx = context.Background()
	}

	token := conn.Query("token", "")
	device := "web"

	var clientId, nickName, avatar string
	if token != "" {
		loginId := auth.CheckToken(cctx, device, token)
		if loginId != "" {
			clientId = loginId
			if v, err := auth.GetSessField(cctx, loginId, "nickName"); err == nil && v != nil {
				if s, ok := v.(string); ok {
					nickName = s
				}
			}
			if nickName == "" {
				nickName = "用户" + loginId
			}
			if v, err := auth.GetSessField(cctx, loginId, "avatar"); err == nil && v != nil {
				if s, ok := v.(string); ok {
					avatar = conf.FileUrl(s)
				}
			}
		}
	}
	if clientId == "" {
		clientId = "0"
		nickName = "游客"
	}

	connId := strconv.FormatInt(idgen.NextId(), 10)
	clientId = clientId + "_" + connId

	client := &ImClient{
		ID:        clientId,
		NickName:  nickName,
		WriteChan: make(chan []byte, writeChanSize),
	}

	var joinedRooms []string
	var roomsMu sync.Mutex

	defer func() {
		// 连接断开时，向各房间广播下线通知（先移除再广播，使 memberCount 正确），再关闭
		close(client.WriteChan)
		roomsMu.Lock()
		rooms := make([]string, len(joinedRooms))
		copy(rooms, joinedRooms)
		roomsMu.Unlock()
		LeaveAllRooms(clientId, rooms)
		for _, rid := range rooms {
			broadcastOffline(rid, clientId, nickName, avatar, GetRoomMemberCount(rid))
		}
		slog.Debug("sys_im connection closed", "clientId", clientId)
	}()

	// 写协程
	go func() {
		for data := range client.WriteChan {
			if err := conn.WriteMessage(fiberws.TextMessage, data); err != nil {
				return
			}
		}
	}()

	// 读循环
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var cm clientMsg
		if err := json.Unmarshal(msg, &cm); err != nil {
			sendError(client, "invalid json")
			continue
		}

		switch cm.Type {
		case "join":
			if cm.RoomId == "" {
				sendError(client, "roomId required")
				continue
			}
			JoinRoom(cm.RoomId, clientId, client)
			roomsMu.Lock()
			joinedRooms = append(joinedRooms, cm.RoomId)
			roomsMu.Unlock()
			count := GetRoomMemberCount(cm.RoomId)
			sendJoined(client, cm.RoomId, count)
			// 向房间其他成员广播上线通知
			broadcastOnline(cm.RoomId, clientId, nickName, avatar, count)

		case "leave":
			if cm.RoomId == "" {
				sendError(client, "roomId required")
				continue
			}
			// 先广播下线通知（此时仍在房间内），再移除
			broadcastOffline(cm.RoomId, clientId, nickName, avatar, GetRoomMemberCount(cm.RoomId)-1)
			LeaveRoom(cm.RoomId, clientId)
			roomsMu.Lock()
			for i, r := range joinedRooms {
				if r == cm.RoomId {
					joinedRooms = append(joinedRooms[:i], joinedRooms[i+1:]...)
					break
				}
			}
			roomsMu.Unlock()
			sendLeft(client, cm.RoomId, GetRoomMemberCount(cm.RoomId))

		case "message":
			if cm.RoomId == "" || cm.Content == "" {
				sendError(client, "roomId and content required")
				continue
			}
			sm := buildServerMsg("message", cm.RoomId, clientId, nickName, avatar, cm.Content, "", time.Now().UnixMilli(), 0)
			data, _ := json.Marshal(sm)
			BroadcastToRoom(cm.RoomId, data, clientId)
			// 回显给发送者，保证本人消息也有完整的 from/nickName/avatar，避免前端用 info 接口数据拼装导致昵称不全
			select {
			case client.WriteChan <- data:
			default:
			}

		default:
			sendError(client, "unknown type: "+cm.Type)
		}
	}
}

func sendError(client *ImClient, msg string) {
	data, _ := json.Marshal(serverMsg{Type: "error", Msg: msg})
	select {
	case client.WriteChan <- data:
	default:
	}
}

func sendJoined(client *ImClient, roomId string, memberCount int) {
	data, _ := json.Marshal(serverMsg{Type: "joined", RoomId: roomId, MemberCount: memberCount})
	select {
	case client.WriteChan <- data:
	default:
	}
}

func sendLeft(client *ImClient, roomId string, memberCount int) {
	data, _ := json.Marshal(serverMsg{Type: "left", RoomId: roomId, MemberCount: memberCount})
	select {
	case client.WriteChan <- data:
	default:
	}
}

// buildServerMsg 构建服务端消息，从 from(loginId_connId) 解析出 loginId、connId
// memberCount 为 0 时 omitempty 不序列化；非 0 时推送房间在线人数
func buildServerMsg(typ, roomId, from, nickName, avatar, content, msg string, t int64, memberCount int) serverMsg {
	loginId, connId := parseFrom(from)
	return serverMsg{
		Type:        typ,
		RoomId:      roomId,
		From:        from,
		LoginId:     loginId,
		ConnId:      connId,
		NickName:    nickName,
		Avatar:      avatar,
		Content:     content,
		Msg:         msg,
		Time:        t,
		MemberCount: memberCount,
	}
}

func parseFrom(from string) (loginId, connId string) {
	if idx := strings.Index(from, "_"); idx >= 0 {
		return from[:idx], from[idx+1:]
	}
	return from, ""
}

// broadcastOnline 向房间内其他成员广播上线通知（排除上线者本人）
func broadcastOnline(roomId, from, nickName, avatar string, memberCount int) {
	sm := buildServerMsg("online", roomId, from, nickName, avatar, "", nickName+" 上线了", time.Now().UnixMilli(), memberCount)
	data, _ := json.Marshal(sm)
	BroadcastToRoom(roomId, data, from)
}

// broadcastOffline 向房间内其他成员广播下线通知（排除下线者本人），memberCount 为移除后的在线人数
func broadcastOffline(roomId, from, nickName, avatar string, memberCount int) {
	sm := buildServerMsg("offline", roomId, from, nickName, avatar, "", nickName+" 下线了", time.Now().UnixMilli(), memberCount)
	data, _ := json.Marshal(sm)
	BroadcastToRoom(roomId, data, from)
}
