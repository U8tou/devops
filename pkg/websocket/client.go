package websocket

import (
	"fmt"
	"sync"
)

var clientMap sync.Map

// Client WebSocket 客户端，支持全双工通信
type Client struct {
	ID        string
	WriteChan chan []byte
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetWriteChan() chan []byte {
	return c.WriteChan
}

func (c *Client) GetClient(id string) *Client {
	v, ok := clientMap.Load(id)
	if !ok {
		return nil
	}
	return v.(*Client)
}

func (c *Client) SetClient(id string, client *Client) {
	clientMap.Store(id, client)
}

func (c *Client) DelClient(id string) {
	clientMap.Delete(id)
}

// SendTo 向指定客户端非阻塞发送，客户端不存在或 channel 满时返回 false
func SendTo(id string, data []byte) bool {
	var client Client
	c := client.GetClient(id)
	if c == nil {
		return false
	}
	select {
	case c.WriteChan <- data:
		return true
	default:
		return false
	}
}

// SendToText 向指定客户端发送文本消息
func SendToText(id string, format string, args ...any) bool {
	return SendTo(id, fmt.Appendf(nil, format, args...))
}
