package alog

import (
	"fmt"
	"log/slog"
	"runtime"
)

func Debug(tag string, data any) {
	if enable {
		slog.Debug(tag,
			slog.Any("data", data),
			slog.Any("source", getSource(2)),
		)
		return
	}
	slog.Debug(tag, slog.Any("data", data))
}
func Info(tag string, data any) {
	if enable {
		slog.Info(tag,
			slog.Any("data", data),
			slog.Any("source", getSource(2)),
		)
		return
	}
	slog.Info(tag, slog.Any("data", data))
}

func Warn(tag string, data any) {
	if enable {
		slog.Warn(tag,
			slog.Any("data", data),
			slog.Any("source", getSource(2)),
		)
		return
	}
	slog.Warn(tag, slog.Any("data", data))
}

func Error(tag string, data any) {
	if enable {
		slog.Error(tag,
			slog.Any("data", data),
			slog.Any("source", getSource(2)),
		)
		return
	}
	slog.Error(tag, slog.Any("data", data))
}

// getSource 返回调用方信息（文件 + 行号 + 函数名）
func getSource(skip int) Source {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	fnName := runtime.FuncForPC(pc).Name()
	return Source{
		File:     fmt.Sprintf("%s:%d", file, line),
		Function: fnName,
	}
}
