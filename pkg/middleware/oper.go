package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var sensitiveRoutes = map[string]bool{
	"/api/sys_auth/login":           true,
	"/api/sys_auth/change_password": true,
	"/api/sys_user/reset_pwd":       true,
	"/api/sys_user/import":          true,
}

// Oper 操作记录中间件
func Oper(log ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		elapsed := time.Since(start)

		method := c.Method()
		path := c.Path()

		fmt.Printf("[%s] %s %s %d %s\n", time.Now().Format("15:04:05"), method, path, c.Response().StatusCode(), elapsed)

		if !sensitiveRoutes[strings.ToLower(path)] {
			if method == "GET" {
				fmt.Printf("  Query: %s\n", c.Context().QueryArgs())
			} else if strings.Contains(c.Get("Content-Type"), "application/json") {
				fmt.Printf("  Body: %s\n", string(c.Body()))
			}
		}

		if err != nil {
			fmt.Printf("  Error: %s\n", err.Error())
		}
		return err
	}
}
