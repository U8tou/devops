package alog

import (
	"log/slog"
	"strings"
)

type Source struct {
	File     string `json:"file"`
	Function string `json:"function"`
}

var (
	opt    Opt  // 全局日志配置
	enable bool // 是否附加调用方 Source 信息
)

// Opt 日志初始化选项
type Opt struct {
	Level      string
	StdOut     bool
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

const (
	defaultLogLevel      = "debug"
	defaultLogFilename   = "logs/app.log"
	defaultLogMaxSize    = 10 // MB
	defaultLogMaxAge     = 28 // 天
	defaultLogMaxBackups = 3
)

// NewOpt 创建空的日志选项，由调用方显式填充配置；未填项在 Build 时按默认值补全
func NewOpt() *Opt {
	return &Opt{}
}

// optIsUnset 表示调用方未做任何配置（全零值），此时布尔项按默认 true 处理。
// 注意：无法区分 YAML/配置里「省略」与「显式 false」，全零以外的 Opt 会保留传入的布尔值。
func optIsUnset(o Opt) bool {
	return o.Level == "" && !o.StdOut && o.Filename == "" &&
		o.MaxSize == 0 && o.MaxAge == 0 && o.MaxBackups == 0 && !o.LocalTime && !o.Compress
}

func (o *Opt) applyFieldDefaults() {
	if o.Level == "" {
		o.Level = defaultLogLevel
	}
	if o.Filename == "" {
		o.Filename = defaultLogFilename
	}
	if o.MaxSize == 0 {
		o.MaxSize = defaultLogMaxSize
	}
	if o.MaxAge == 0 {
		o.MaxAge = defaultLogMaxAge
	}
	if o.MaxBackups == 0 {
		o.MaxBackups = defaultLogMaxBackups
	}
}

// Build 根据选项初始化全局 slog.Logger；零值字段会补全为上述默认配置
func (o *Opt) Build() {
	raw := *o
	fullUnset := optIsUnset(raw)
	raw.applyFieldDefaults()
	if fullUnset {
		raw.StdOut = true
		raw.Compress = true
		raw.LocalTime = true
	}
	opt = raw

	// 日志级别为 debug 时才附加 Source 信息
	enable = strings.EqualFold(opt.Level, "debug")
	slog.SetDefault(newLogger())
}
