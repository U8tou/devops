package middleware

import (
	"pkg/auth"
	"pkg/conf"
	"pkg/constx"

	"github.com/gofiber/fiber/v2"
)

// Auth 认证处理中间件
// permis: 可选的权限列表，如果提供则检查用户是否拥有任意一个权限
func Auth(permis ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 获取 Token（Header → 与 tokenName 同名的 query → `token` query，与 sys_sse / EventSource 一致）
		token := c.Get(conf.Auth.TokenName)
		if token == "" {
			token = c.Query(conf.Auth.TokenName)
		}
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			return fiber.ErrUnauthorized
		}

		// 获取设备标识（默认为 web）
		device := c.Get("X-Device", "web")

		// 校验 Token
		loginId := auth.CheckToken(c.Context(), device, token)
		if loginId == "" {
			return fiber.ErrUnauthorized
		}

		// 设置登录用户ID到上下文
		constx.SetLoginId(c, loginId)

		// 如果需要检查权限（根用户 userType=00 跳过）
		if len(permis) > 0 && !auth.IsRootUser(c.Context(), loginId) {
			if !auth.HasAnyPermis(c.Context(), loginId, permis) {
				return fiber.ErrForbidden
			}
		}

		return c.Next()
	}
}
