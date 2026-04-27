package sse

import "strings"

// SendTo 向指定客户端非阻塞发送消息，客户端不存在或 channel 满时返回 false
func SendTo(sseId string, data []byte) bool {
	var client SseClient
	c := client.GetCli(sseId)
	if c == nil {
		return false
	}
	select {
	case c.writeChan <- data:
		return true
	default:
		return false
	}
}

// TargetedText 向指定 SSE 客户端发送文本消息（非阻塞）。
// tag 建议使用常量：TagGreeting/TagMessage/TagNotify/TagWarning；msg 为消息正文。
func TargetedText(sseId, tag, msg string) bool {
	return SendTo(sseId, sseMessage(tag, msg))
}

// SendToUser 向某用户的所有 SSE 连接发送消息（按 username_ 前缀匹配，如 1_*）。
// 返回成功送达的连接数。用于通知指定登录用户，该用户多端在线时都能收到。
func SendToUser(username, tag, msg string) int {
	prefix := username + "_"
	data := sseMessage(tag, msg)
	var n int
	SseMap.Range(func(k, v any) bool {
		if strings.HasPrefix(k.(string), prefix) {
			c := v.(*SseClient)
			select {
			case c.writeChan <- data:
				n++
			default:
			}
		}
		return true
	})
	return n
}
