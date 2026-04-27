package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"pkg/alog"
	"pkg/auth"
	"pkg/conf"
	"pkg/db"
	"pkg/middleware"
	"pkg/rds"
	"pkg/resp"
	"time"

	"admin/docs"
	devprocess "admin/internal/dev_process"
	sysim "admin/internal/sys_im"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/yitter/idgenerator-go/idgen"
)

func Start() {
	// 基础能力初始化
	appConf := conf.App
	initID(appConf.Id)
	initAuth(conf.Auth)
	initLog(conf.Log)
	initDb()
	devprocess.StartCronScheduler()
	sysim.InitDefaultRoom()

	// 启动服务
	app := fiber.New(fiber.Config{
		AppName:       appConf.Name,
		ServerHeader:  appConf.Doc,
		StrictRouting: true, // /foo 和 /foo/ 被视为不同路由
		CaseSensitive: true, // /Foo 和 /foo 被视为不同路由
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 0, // SSE/WebSocket 长连接必须为 0，fasthttp 的 WriteTimeout 从请求开始计费，心跳无法续期
		BodyLimit:     50 * 1024 * 1024, // 50M 响应体大小
		JSONEncoder:   sonic.Marshal,    // 换成更快的 JSON 库
		JSONDecoder:   sonic.Unmarshal,
		// Prefork:       true, // 启动多个进程来处理连接
		// EnablePrintRoutes: true, // 启动时打印路由
		// DisableStartupMessage: true, // 是否不打印Banner
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			alog.Error("❌ 异常:", err.Error())
			return resp.Error(ctx, err)
		},
	})
	// 恐慌恢复
	app.Use(recover.New())
	// CORS for external resources
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Shh-Encrypted",
	}))
	// 設置安全請求頭
	app.Use(helmet.New())
	// 并发限制
	app.Use(limiter.New(limiter.Config{
		Max:        200, // 每 IP 并发上限
		Expiration: 2 * time.Second,
	}))
	// 压缩（跳过 WebSocket/SSE 长连接，避免握手异常）
	app.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
		Next: func(c *fiber.Ctx) bool {
			path := c.Path()
			return path == "/api/sys_websocket" || path == "/api/sys_sse" ||
				path == "/api/dev_process/run_stream" ||
				path == "/api/sys_im/chat" || path == "/sys_websocket" || path == "/sys_sse" ||
				path == "/sys_im/chat" // 兼容无 /api 前缀
		},
	}))

	// 網站圖標
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
		URL:  "/favicon.ico",
	}))
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))
	// 静态服务
	app.Static("/", "./dist", fiber.Static{
		Compress:      true,
		ByteRange:     true,
		Browse:        true,
		Index:         "index.html",
		CacheDuration: 24 * time.Hour,
		MaxAge:        3600,
	})
	// 运行命令 swag init -g cmd/main.go -o docs -q --parseDependency --parseInternal --instanceName swagger
	app.Get("/swagger/*", swagger.New(swagger.Config{
		DefaultModelsExpandDepth: -1,   // 隐藏最后的Models
		PersistAuthorization:     true, // 刷新不丢失数据
	}))

	// 接口文档
	docs.SwaggerInfo.Title = appConf.Name
	docs.SwaggerInfo.Description = "接口文档"
	docs.SwaggerInfo.Version = appConf.Version
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", appConf.Port)
	docs.SwaggerInfo.BasePath = appConf.BaseUrl
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// 业务处理
	api := app.Group("/api")
	api.Use(middleware.ApiCrypto())
	// 上传文件静态访问（如头像）
	api.Static("/uploads", conf.File.Local)
	// 注册路由
	registerRoute(api)

	// 404处理
	app.Use(func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})

	// 优雅关闭服务器
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		devprocess.StopCronScheduler()
		if err := app.Shutdown(); err != nil {
			fmt.Printf("关闭服务器错误: %s", err)
		}
	}()
	// 启动服务
	log.Fatal(app.Listen(fmt.Sprintf(":%d", appConf.Port)))
}

// initID 初始化分布式 ID 生成器
func initID(workerID uint16) {
	// 创建 IdGeneratorOptions 对象，可在构造函数中输入 WorkerId：
	options := idgen.NewIdGeneratorOptions(workerID)
	// 如需自定义 BaseTime、WorkerIdBitLength 等，可在此处扩展 options
	idgen.SetIdGenerator(options)
	fmt.Println("✅ ID ok!")
}

// initLog 初始化日志框架（Openobserver 配置见 conf.Openobserver，远端推送接入可后续在 alog 中实现）
func initLog(c conf.LogConf) {
	(&alog.Opt{
		Level:      c.Level,
		StdOut:     c.StdOut,
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxAge:     c.MaxAge,
		MaxBackups: c.MaxBackups,
		LocalTime:  c.LocalTime,
		Compress:   c.Compress,
	}).Build()
	fmt.Println("✅ slog ok!")
}

// initAuth 初始化认证框架
func initAuth(c conf.AuthConf) {
	var aut auth.AuthOpt
	if c.Use == "redis" {
		aut.WithStorage(rds.GetRds())
	} else {
		aut.WithStorage(nil)
	}
	if c.KeyPrefix != "" {
		aut.WithKeyPrefix(c.KeyPrefix)
	}
	if c.RetryCount > 0 {
		aut.WithRetryCount(c.RetryCount)
	}
	if c.RetryTime > 0 {
		aut.WithRetryTime(c.RetryTime)
	}
	if c.Timeout > 0 {
		aut.WithAccTokenTtl(int(c.Timeout))
	}
	aut.Build()
	fmt.Println("✅ auth ok!")
}

func initDb() {
	db.GetDb()
	db.StartBackup()
}
