package syswebsocket

import (
	"encoding/json"
	"log/slog"
	"pkg/errs"
	wsclient "pkg/websocket"

	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const (
	writeChanSize = 100
)

// Upgrader 返回 Fiber WebSocket 处理器
func Upgrader() fiberws.Config {
	return fiberws.Config{
		RecoverHandler: func(conn *fiberws.Conn) {
			if err := recover(); err != nil {
				slog.Error("WebSocket panic", "err", err)
			}
		},
	}
}

// @Summary		WebSocket 全双工连接
// @Description	建立 WebSocket 连接，支持服务端推送与客户端上报；需传 uuid 标识连接，通过 syswebsocket.SendMsg(uuid, tag, msg) 推送；客户端可发送文本/二进制，服务端会 echo 文本消息
// @Tags			SysWebSocketApi
// @Accept			json
// @Produce		json
// @Param			uuid	query		string	true	"连接标识，用于后续 SendMsg 定向推送"
// @Success		101		{string}	string	"Switching Protocols，升级为 WebSocket"
// @Failure		400		{string}	string	"uuid 为空或非 WebSocket 升级请求"
// @Router			/sys_websocket [get]
func Handler(conn *fiberws.Conn) {
	uuid := conn.Query("uuid", "")
	if uuid == "" {
		_ = conn.WriteJSON(fiber.Map{"error": errs.New("uuid is required").Error()})
		_ = conn.Close()
		return
	}

	client := &wsclient.Client{
		ID:        uuid,
		WriteChan: make(chan []byte, writeChanSize),
	}
	client.SetClient(uuid, client)

	defer func() {
		close(client.WriteChan)
		client.DelClient(uuid)
		slog.Debug("WebSocket connection closed", "uuid", uuid)
	}()

	// 写协程：将 channel 中的消息发送给客户端
	go func() {
		for data := range client.WriteChan {
			if err := conn.WriteMessage(fiberws.TextMessage, data); err != nil {
				return
			}
		}
	}()

	// 读循环：接收客户端消息（全双工）
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		// 可根据 mt 和 msg 做业务处理，例如 echo、路由等
		switch mt {
		case fiberws.TextMessage:
			slog.Debug("WebSocket recv text", "uuid", uuid, "msg", string(msg))
			// 示例：可在此解析消息并响应
			_ = conn.WriteMessage(fiberws.TextMessage, msg)
		case fiberws.BinaryMessage:
			slog.Debug("WebSocket recv binary", "uuid", uuid)
		case fiberws.PingMessage, fiberws.PongMessage:
			// 由底层自动处理
		case fiberws.CloseMessage:
			return
		}
	}
}

// SendMsg 向指定 WebSocket 客户端发送消息，非阻塞
func SendMsg(uuid string, tag string, msg string) {
	payload, _ := json.Marshal(fiber.Map{
		"event": "message",
		"tag":   tag,
		"data":  msg,
	})
	wsclient.SendTo(uuid, payload)
}
