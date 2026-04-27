package sse

import (
	"fmt"
	"sync"
)

var SseMap sync.Map

// 消息类型标签
const (
	TagGreeting = "greeting" // 问候语
	TagMessage  = "message"  // 消息
	TagNotify   = "notify"   // 通知
	TagWarning  = "warning"  // 警告
)

type SseClient struct {
	sseId     string
	writeChan chan []byte
}

func (m *SseClient) Create(sseId string) {
	m.sseId = sseId
	m.writeChan = make(chan []byte, 100)
}

func (m *SseClient) GetSseId() string {
	return m.sseId
}

func (m *SseClient) GetData() chan []byte {
	return m.writeChan
}

func (m *SseClient) GetCli(sseId string) *SseClient {
	v, ok := SseMap.Load(sseId)
	if !ok {
		return nil
	}
	return v.(*SseClient)
}

func (m *SseClient) SetCli(sseId string, client *SseClient) {
	SseMap.Store(sseId, client)
}

func (m *SseClient) DelCli(sseId string) {
	SseMap.Delete(sseId)
}

// Send 非阻塞发送，channel 满或客户端不存在时丢弃消息，避免阻塞调用方
func (m *SseClient) Send(data []byte) bool {
	client := m.GetCli(m.sseId)
	if client == nil {
		return false
	}
	select {
	case client.writeChan <- data:
		return true
	default:
		return false
	}
}

// Sendf 格式化并非阻塞发送
func (m *SseClient) Sendf(format string, args ...any) bool {
	return m.Send(fmt.Appendf(nil, format, args...))
}

func sseMessage(tag, msg string) []byte {
	return fmt.Appendf(nil, "event: message\ndata:%s-%s\n\n", tag, msg)
}
