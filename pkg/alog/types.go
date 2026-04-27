package alog

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/fatih/color"
	"github.com/tidwall/gjson"
)

type CustomHandler struct {
	slog.Handler
}

func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	// 将时间戳转为指定格式（如 Unix 时间戳或 ISO8601）
	// r.AddAttrs(slog.String("timestamp", r.Time.Format(time.RFC3339)))
	// 或转为 Unix 时间戳（整数）
	r.AddAttrs(slog.Int64("_timestamp", r.Time.UnixMicro()))
	return h.Handler.Handle(ctx, r) // 调用底层 Handler
}

// 自定义实现io.Writer接口，将数据输出到自定义位置
type myWriter struct {
	ToBuffer bool
	ToStdOut bool
}

// Write 实现io.Writer接口的Write方法
func (rw *myWriter) Write(p []byte) (n int, err error) {
	// 是否输出到控制台
	if rw.ToStdOut {
		lg := string(p)
		// 控制台打印
		res := gjson.GetMany(lg, "data", "level", "time", "msg", "source.file")
		// 输出控制器
		var logFunc func(format string, a ...any)
		switch res[1].Str {
		case "INFO":
			logFunc = color.Blue
		case "WARN":
			logFunc = color.Yellow
		case "ERROR":
			logFunc = color.Red
		default:
			logFunc = color.Cyan
		}
		if res[0].Exists() {
			logFunc("[alog] [%s]  %s %s [%s]\n - %s\n",
				res[1].Str,
				res[2].Str,
				res[3].Str,
				res[4].Str,
				res[0].Raw,
			)
		} else {
			logFunc("[alog] [%s]  %s %s [%s]\n",
				res[1].Str,
				res[2].Str,
				res[3].Str,
				res[4].Str,
			)
		}
	}
	return len(p), nil
}

func BytesSliceToJSONArray(byteSlices [][]byte) ([]byte, error) {
	var result []any
	for _, bytes := range byteSlices {
		var data any
		if err := json.Unmarshal(bytes, &data); err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return json.Marshal(result)
}
