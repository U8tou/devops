package sse

// Broadcast 向所有在线客户端广播原始数据（非阻塞），channel 满的客户端会丢弃
func Broadcast(data []byte) int {
	var n int
	SseMap.Range(func(_, v any) bool {
		c := v.(*SseClient)
		select {
		case c.writeChan <- data:
			n++
		default:
		}
		return true
	})
	return n
}

// BroadcastText 向所有在线客户端广播文本消息（非阻塞）。
// tag 建议使用常量：TagGreeting/TagMessage/TagNotify/TagWarning；msg 为消息正文。
func BroadcastText(tag, msg string) int {
	return Broadcast(sseMessage(tag, msg))
}
