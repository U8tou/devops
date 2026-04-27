package sysim

import (
	"sync"
)

const defaultRoomId = "1"

var (
	roomMap sync.Map // roomId -> *roomClients
)

type roomClients struct {
	mu      sync.RWMutex
	clients map[string]*ImClient // clientId -> client
}

func getOrCreateRoom(roomId string) *roomClients {
	v, _ := roomMap.LoadOrStore(roomId, &roomClients{
		clients: make(map[string]*ImClient),
	})
	return v.(*roomClients)
}

// InitDefaultRoom 系统启动时初始化默认聊天室（ID 为 1）
func InitDefaultRoom() {
	getOrCreateRoom(defaultRoomId)
}

// JoinRoom 加入房间
func JoinRoom(roomId, clientId string, client *ImClient) {
	if roomId == "" {
		return
	}
	rc := getOrCreateRoom(roomId)
	rc.mu.Lock()
	rc.clients[clientId] = client
	rc.mu.Unlock()
}

// GetRoomMemberCount 获取房间当前在线人数
func GetRoomMemberCount(roomId string) int {
	if roomId == "" {
		return 0
	}
	v, ok := roomMap.Load(roomId)
	if !ok {
		return 0
	}
	rc := v.(*roomClients)
	rc.mu.RLock()
	n := len(rc.clients)
	rc.mu.RUnlock()
	return n
}

// LeaveRoom 离开房间
func LeaveRoom(roomId, clientId string) {
	if roomId == "" {
		return
	}
	rc := getOrCreateRoom(roomId)
	rc.mu.Lock()
	delete(rc.clients, clientId)
	rc.mu.Unlock()
}

// LeaveAllRooms 从所有房间移除该客户端（连接断开时调用）
func LeaveAllRooms(clientId string, roomIds []string) {
	for _, rid := range roomIds {
		LeaveRoom(rid, clientId)
	}
}

// BroadcastToRoom 向房间内所有成员广播（excludeId 排除发送者）
func BroadcastToRoom(roomId string, data []byte, excludeId string) int {
	if roomId == "" {
		return 0
	}
	v, ok := roomMap.Load(roomId)
	if !ok {
		return 0
	}
	rc := v.(*roomClients)
	rc.mu.RLock()
	clients := make([]*ImClient, 0, len(rc.clients))
	for id, c := range rc.clients {
		if id != excludeId {
			clients = append(clients, c)
		}
	}
	rc.mu.RUnlock()

	var n int
	for _, c := range clients {
		select {
		case c.WriteChan <- data:
			n++
		default:
		}
	}
	return n
}
