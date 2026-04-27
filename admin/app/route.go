package app

import (
	devprocess "admin/internal/dev_process"
	devproject "admin/internal/dev_project"
	sysauth "admin/internal/sys_auth"
	syscommon "admin/internal/sys_common"
	sysdept "admin/internal/sys_dept"
	sysim "admin/internal/sys_im"
	sysmenu "admin/internal/sys_menu"
	syspost "admin/internal/sys_post"
	sysrole "admin/internal/sys_role"
	syssse "admin/internal/sys_sse"
	sysuser "admin/internal/sys_user"
	syswebsocket "admin/internal/sys_websocket"
	"errors"
	m "pkg/middleware"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	t "github.com/gofiber/fiber/v2/middleware/timeout"
)

var (
	ErrRequestTimeout = errors.New("request context canceled or timeout")
	TimeOut           = 2 * time.Second
)

func registerRoute(r fiber.Router) {
	// SSE（服务端推送，长连接不使用 timeout）
	r.Get("/sys_sse", syssse.Sse)
	// WebSocket（全双工）
	r.Use("/sys_websocket", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	r.Get("/sys_websocket", websocket.New(syswebsocket.Handler, syswebsocket.Upgrader()))
	// IM 聊天室
	r.Use("/sys_im", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("ctx", c.Context())
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	r.Get("/sys_im/chat", websocket.New(sysim.Handler, sysim.Upgrader()))
	// 公共模块
	r.Route("/sys_common", func(api fiber.Router) {
		api.Get("/get_file/:dir/:obj", m.Oper(), t.NewWithContext(syscommon.GetFile, 10*time.Second, ErrRequestTimeout))
	})
	// 认证模块
	r.Route("/sys_auth", func(api fiber.Router) {
		api.Post("/login", m.Oper(), t.NewWithContext(sysauth.Login, TimeOut, ErrRequestTimeout))
		api.Post("/register", m.Oper(), t.NewWithContext(sysauth.Register, TimeOut, ErrRequestTimeout))
		api.Post("/refresh_token", m.Oper(), t.NewWithContext(sysauth.RefreshToken, TimeOut, ErrRequestTimeout))
		api.Get("/info", m.Auth(), m.Oper(), t.NewWithContext(sysauth.Info, TimeOut))
		api.Post("/logout", m.Auth(), t.NewWithContext(sysauth.Logout, TimeOut))
		api.Post("/change_password", m.Auth(), t.NewWithContext(sysauth.ChangePassword, TimeOut))
		api.Post("/profile", m.Auth(), t.NewWithContext(sysauth.UpdateProfile, TimeOut))
		api.Post("/avatar", m.Auth(), t.NewWithContext(sysauth.UploadAvatar, TimeOut))
	})
	// 用户管理
	r.Route("/sys_user", func(api fiber.Router) {
		api.Get("/list", m.Auth("sys:user:list"), m.Oper(), t.NewWithContext(sysuser.List, TimeOut))
		api.Get("/get", m.Auth("sys:user:get"), m.Oper(), t.NewWithContext(sysuser.Get, TimeOut))
		api.Delete("/del", m.Auth("sys:user:del"), m.Oper(), t.NewWithContext(sysuser.Del, TimeOut))
		api.Post("/add", m.Auth("sys:user:add"), m.Oper(), t.NewWithContext(sysuser.Add, TimeOut))
		api.Put("/edit", m.Auth("sys:user:edit"), m.Oper(), t.NewWithContext(sysuser.Edit, TimeOut))
		api.Post("/reset_pwd", m.Auth("sys:user:reset_pwd"), m.Oper(), t.NewWithContext(sysuser.ResetPassword, TimeOut))
		api.Get("/all", m.Auth("sys:user:get"), m.Oper(), t.NewWithContext(sysuser.All, TimeOut))
		api.Post("/assign_role", m.Auth("sys:user:assign_role"), m.Oper(), t.NewWithContext(sysuser.AssignRole, TimeOut))
		api.Post("/assign_dept", m.Auth("sys:user:assign_dept"), m.Oper(), t.NewWithContext(sysuser.AssignDept, TimeOut))
		api.Post("/assign_post", m.Auth("sys:user:assign_post"), m.Oper(), t.NewWithContext(sysuser.AssignPost, TimeOut))
		api.Get("/template", m.Auth("sys:user:list"), m.Oper(), t.NewWithContext(sysuser.Template, TimeOut, ErrRequestTimeout))
		api.Get("/export", m.Auth("sys:user:list"), m.Oper(), t.NewWithContext(sysuser.Export, 10*time.Second, ErrRequestTimeout))
		api.Post("/import", m.Auth("sys:user:add"), m.Oper(), t.NewWithContext(sysuser.Import, 30*time.Second, ErrRequestTimeout))
	})
	// 角色管理
	r.Route("/sys_role", func(api fiber.Router) {
		api.Get("/list", m.Auth("sys:role:list"), m.Oper(), t.NewWithContext(sysrole.List, TimeOut))
		api.Get("/get", m.Auth("sys:role:get"), m.Oper(), t.NewWithContext(sysrole.Get, TimeOut))
		api.Delete("/del", m.Auth("sys:role:del"), m.Oper(), t.NewWithContext(sysrole.Del, TimeOut))
		api.Post("/add", m.Auth("sys:role:add"), m.Oper(), t.NewWithContext(sysrole.Add, TimeOut))
		api.Put("/edit", m.Auth("sys:role:edit"), m.Oper(), t.NewWithContext(sysrole.Edit, TimeOut))
		api.Get("/all", m.Auth(), m.Oper(), t.NewWithContext(sysrole.All, TimeOut))
		api.Post("/assign_dept", m.Auth("sys:role:assign_dept"), m.Oper(), t.NewWithContext(sysrole.AssignDept, TimeOut))
		api.Post("/assign_menu", m.Auth("sys:role:assign_menu"), m.Oper(), t.NewWithContext(sysrole.AssignMenu, TimeOut))
	})
	// 部门管理
	r.Route("/sys_dept", func(api fiber.Router) {
		api.Get("/list", m.Auth("sys:dept:list"), m.Oper(), t.NewWithContext(sysdept.List, TimeOut))
		api.Get("/get", m.Auth("sys:dept:get"), m.Oper(), t.NewWithContext(sysdept.Get, TimeOut))
		api.Delete("/del", m.Auth("sys:dept:del"), m.Oper(), t.NewWithContext(sysdept.Del, TimeOut))
		api.Post("/add", m.Auth("sys:dept:add"), m.Oper(), t.NewWithContext(sysdept.Add, TimeOut))
		api.Put("/edit", m.Auth("sys:dept:edit"), m.Oper(), t.NewWithContext(sysdept.Edit, TimeOut))
		api.Get("/all", m.Auth(), m.Oper(), t.NewWithContext(sysdept.All, TimeOut))
	})
	// 岗位管理
	r.Route("/sys_post", func(api fiber.Router) {
		api.Get("/list", m.Auth("sys:post:list"), m.Oper(), t.NewWithContext(syspost.List, TimeOut))
		api.Get("/get", m.Auth("sys:post:get"), m.Oper(), t.NewWithContext(syspost.Get, TimeOut))
		api.Delete("/del", m.Auth("sys:post:del"), m.Oper(), t.NewWithContext(syspost.Del, TimeOut))
		api.Post("/add", m.Auth("sys:post:add"), m.Oper(), t.NewWithContext(syspost.Add, TimeOut))
		api.Put("/edit", m.Auth("sys:post:edit"), m.Oper(), t.NewWithContext(syspost.Edit, TimeOut))
		api.Get("/all", m.Auth(), m.Oper(), t.NewWithContext(syspost.All, TimeOut))
	})
	// 菜单/权限管理
	r.Route("/sys_menu", func(api fiber.Router) {
		api.Get("/list", m.Auth("sys:menu:list"), m.Oper(), t.NewWithContext(sysmenu.List, TimeOut))
		api.Get("/get", m.Auth("sys:menu:get"), m.Oper(), t.NewWithContext(sysmenu.Get, TimeOut))
		api.Get("/all", m.Auth(), m.Oper(), t.NewWithContext(sysmenu.All, TimeOut))
	})
	// 自动化流程
	r.Route("/dev_process", func(api fiber.Router) {
		api.Get("/list", m.Auth("dev:process:list"), m.Oper(), t.NewWithContext(devprocess.List, TimeOut))
		api.Get("/get", m.Auth("dev:process:get"), m.Oper(), t.NewWithContext(devprocess.Get, TimeOut))
		api.Post("/add", m.Auth("dev:process:add"), m.Oper(), t.NewWithContext(devprocess.Add, TimeOut))
		api.Put("/edit", m.Auth("dev:process:edit"), m.Oper(), t.NewWithContext(devprocess.Edit, TimeOut))
		api.Put("/edit_flow", m.Auth("dev:process:edit"), m.Oper(), t.NewWithContext(devprocess.EditFlow, TimeOut))
		api.Put("/edit_env", m.Auth("dev:process:edit"), m.Oper(), t.NewWithContext(devprocess.EditEnv, TimeOut))
		api.Put("/set_cron_enabled", m.Auth("dev:process:edit"), m.Oper(), t.NewWithContext(devprocess.SetCronEnabled, TimeOut))
		// 录入时校验节点：与编辑流程同属配置阶段，允许 edit 或 run 任一权限（避免仅有编辑权却 403）
		api.Post("/validate_node", m.Auth("dev:process:edit", "dev:process:run"), m.Oper(), t.NewWithContext(devprocess.ValidateNode, 30*time.Second, ErrRequestTimeout))
		// 执行/流式执行/取消：与编辑流程一致，具备 edit 或 run 任一即可（避免仅有编辑权无法点「执行」）
		api.Get("/run_stream", m.Auth("dev:process:edit", "dev:process:run"), m.Oper(), devprocess.RunStream)
		api.Post("/run_cancel", m.Auth("dev:process:edit", "dev:process:run"), m.Oper(), t.NewWithContext(devprocess.RunCancel, TimeOut))
		api.Post("/run", m.Auth("dev:process:edit", "dev:process:run"), m.Oper(), t.NewWithContext(devprocess.Run, 35*time.Minute, ErrRequestTimeout))
		api.Delete("/del", m.Auth("dev:process:del"), m.Oper(), t.NewWithContext(devprocess.Del, TimeOut))
		// 流程标签字典（新标签权限与原有 list/edit 并存，兼容旧角色）
		api.Get("/tag/list", m.Auth("dev:process:list", "dev:process:tag:list", "dev:process:tag:add", "dev:process:tag:edit", "dev:process:tag:del"), m.Oper(), t.NewWithContext(devprocess.TagList, TimeOut))
		api.Post("/tag/add", m.Auth("dev:process:edit", "dev:process:tag:add"), m.Oper(), t.NewWithContext(devprocess.TagAdd, TimeOut))
		api.Put("/tag/edit", m.Auth("dev:process:edit", "dev:process:tag:edit"), m.Oper(), t.NewWithContext(devprocess.TagEdit, TimeOut))
		api.Delete("/tag/del", m.Auth("dev:process:edit", "dev:process:tag:del"), m.Oper(), t.NewWithContext(devprocess.TagDel, TimeOut))
	})
	// 项目管理
	r.Route("/dev_project", func(api fiber.Router) {
		api.Get("/list", m.Auth("dev:project:list"), m.Oper(), t.NewWithContext(devproject.List, TimeOut))
		api.Get("/get", m.Auth("dev:project:get"), m.Oper(), t.NewWithContext(devproject.Get, TimeOut))
		api.Post("/add", m.Auth("dev:project:add"), m.Oper(), t.NewWithContext(devproject.Add, TimeOut))
		api.Put("/edit", m.Auth("dev:project:edit"), m.Oper(), t.NewWithContext(devproject.Edit, TimeOut))
		api.Put("/edit_mind", m.Auth("dev:project:edit"), m.Oper(), t.NewWithContext(devproject.EditMind, 30*time.Second, ErrRequestTimeout))
		api.Delete("/del", m.Auth("dev:project:del"), m.Oper(), t.NewWithContext(devproject.Del, TimeOut))
		api.Get("/tag/list", m.Auth("dev:project:list", "dev:project:tag:list", "dev:project:tag:add", "dev:project:tag:edit", "dev:project:tag:del"), m.Oper(), t.NewWithContext(devproject.TagList, TimeOut))
		api.Post("/tag/add", m.Auth("dev:project:edit", "dev:project:tag:add"), m.Oper(), t.NewWithContext(devproject.TagAdd, TimeOut))
		api.Put("/tag/edit", m.Auth("dev:project:edit", "dev:project:tag:edit"), m.Oper(), t.NewWithContext(devproject.TagEdit, TimeOut))
		api.Delete("/tag/del", m.Auth("dev:project:edit", "dev:project:tag:del"), m.Oper(), t.NewWithContext(devproject.TagDel, TimeOut))
	})
}
