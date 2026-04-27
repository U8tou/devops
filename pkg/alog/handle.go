package alog

import (
	"io"
	"log/slog"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// newLogger 根据全局 opt 创建 slog.Logger
func newLogger() *slog.Logger {
	var level slog.Level
	switch strings.ToLower(opt.Level) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}
	confOpt := &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
		// ReplaceAttr: nil,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.AnyValue(a.Value.Time().Format("2006/01/02 15:04:05.000000"))
			}
			return a
		},
	}
	// 存储日志
	return useLocal(confOpt)
}
func useLocal(f *slog.HandlerOptions) *slog.Logger {
	// 使用本地日志
	l := &lumberjack.Logger{
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
		LocalTime:  opt.LocalTime,
	}
	stdout := &myWriter{
		ToBuffer: false,
		ToStdOut: opt.StdOut,
	}
	// 输出到控制台和文件
	jsonHandle := slog.NewJSONHandler(io.MultiWriter(l, stdout), f)
	// 完成配置
	def := slog.New(jsonHandle)
	return def
}
