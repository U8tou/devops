package syssse

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"pkg/auth"
	"pkg/conf"
	"pkg/sse"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/yitter/idgenerator-go/idgen"
)

const (
	heartbeatInterval = 15 * time.Second // 保活间隔
	guestPrefix       = "guest"
)

// 励志问候语
var motivationalQuotes = []string{
	"天道酬勤，每一份努力都会开花结果。",
	"越努力，越幸运，坚持就是胜利。",
	"今日事今日毕，明天会更好。",
	"心怀希望，向阳而生。",
	"行动是治愈恐惧的良药。",
	"做最好的自己，遇见更好的未来。",
	"每一次付出，都是在为成功铺路。",
	"相信自己，你一定可以。",
	"保持热爱，奔赴山海。",
	"星光不问赶路人，时光不负有心人。",
}

// sseComment 生成 SSE 注释格式的保活消息（客户端会忽略）
func sseComment(msg string) []byte {
	return fmt.Appendf(nil, ": %s\n\n", msg)
}

// getTimeGreeting 根据当前时间返回时段问候
func getTimeGreeting() string {
	hour := time.Now().Hour()
	switch {
	case hour >= 0 && hour < 5:
		return "凌晨好"
	case hour >= 5 && hour < 8:
		return "早上好"
	case hour >= 8 && hour < 11:
		return "上午好"
	case hour >= 11 && hour < 13:
		return "中午好"
	case hour >= 13 && hour < 17:
		return "下午好"
	case hour >= 17 && hour < 19:
		return "傍晚好"
	default:
		return "晚上好"
	}
}

// randomMotivational 随机返回一条励志语
func randomMotivational() string {
	return motivationalQuotes[rand.Intn(len(motivationalQuotes))]
}

// buildWelcomeGreeting 构建连接成功的问候语，非游客则展示用户昵称
func buildWelcomeGreeting(ctx context.Context, sseId string) string {
	var name string
	if strings.HasPrefix(sseId, guestPrefix) {
		name = "游客"
	} else {
		// loginId 为 sseId 下划线前的部分，从 session 获取用户昵称
		if idx := strings.Index(sseId, "_"); idx > 0 {
			loginId := sseId[:idx]
			if v, err := auth.GetSessField(ctx, loginId, "nickName"); err == nil && v != nil {
				if s, ok := v.(string); ok && s != "" {
					name = s
				} else {
					name = "用户"
				}
			} else {
				name = "用户"
			}
		} else {
			name = "用户"
		}
	}
	return fmt.Sprintf("%s，%s，%s", getTimeGreeting(), name, randomMotivational())
}

// @Summary		SSE 长连接
// @Description	建立 Server-Sent Events 连接。connId 由后端用分布式 ID 生成，前端无需传；传 token 鉴权：有 token 且有效则为登录用户，sseId=loginId_connId，可被 SendToUser 定向通知；无 token 或无效为游客，sseId=guest_connId，仅能接收广播。
// @Tags			SysSseApi
// @Accept			json
// @Produce		text/event-stream
// @Param			token	query		string	false	"鉴权 token（可选，也可放 Authorization 头）"
// @Success		200		{string}	string	"text/event-stream 流式响应，含 heartbeat 保活"
// @Router			/sys_sse [get]
func Sse(c *fiber.Ctx) error {
	connId := strconv.FormatInt(idgen.NextId(), 10) // 系统分布式 ID

	// token 鉴权：支持 query 或 header（EventSource 常用 query）
	token := c.Query("token")
	if token == "" {
		tokenName := conf.Auth.TokenName
		if tokenName == "" {
			tokenName = "Authorization"
		}
		token = c.Get(tokenName)
	}

	var sseId string
	if token != "" {
		device := c.Get("X-Device", "web")
		loginId := auth.CheckToken(c.Context(), device, token)
		if loginId != "" {
			sseId = loginId + "_" + connId // 登录用户：loginId_connId
		} else {
			sseId = guestPrefix + "_" + connId // token 无效视为游客
		}
	} else {
		sseId = guestPrefix + "_" + connId // 游客：仅能接收广播
	}

	client := &sse.SseClient{}
	client.Create(sseId)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("X-Accel-Buffering", "no") // 禁用 nginx 缓冲

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		client.SetCli(sseId, client)
		// 注册成功后发送问候语（StreamWriter 在独立 goroutine 中运行，c.Context() 可能已失效，使用 context.Background）
		_ = sse.TargetedText(sseId, sse.TagGreeting, buildWelcomeGreeting(context.Background(), sseId))
		ticker := time.NewTicker(heartbeatInterval)
		defer func() {
			ticker.Stop()
			client.DelCli(sseId)
			slog.Debug("SSE connection closed", "sseId", sseId)
		}()

		for {
			select {
			case data, ok := <-client.GetData():
				if !ok {
					return
				}
				if _, err := w.Write(data); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}

			case <-ticker.C:
				if _, err := w.Write(sseComment("heartbeat")); err != nil {
					return
				}
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	}))
	return nil
}

// SendToUser 向指定用户的所有 SSE 连接发送消息（该用户多端在线时都能收到），返回成功送达数
func SendToUser(username, tag, msg string) int {
	return sse.SendToUser(username, tag, msg)
}

// BroadcastText 向所有在线客户端（含登录用户和游客）广播消息
func BroadcastText(tag, msg string) int {
	return sse.BroadcastText(tag, msg)
}
