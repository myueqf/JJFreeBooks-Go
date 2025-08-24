package utils

import (
	"context"
	"io"
	"log/slog"
	"os"
	"sync"
)

// ANSI 颜色码
const (
	colorReset = "\x1b[0m"
	colorDebug = "\x1b[36m" // 青色
	colorInfo  = "\x1b[32m" // 绿色
	colorWarn  = "\x1b[33m" // 黄色
	colorError = "\x1b[31m" // 红色
)

var (
	loggerOnce    sync.Once
	defaultLogger *slog.Logger
)

type ColorTextHandler struct {
	w    io.Writer
	opts *slog.HandlerOptions
}

func NewColorTextHandler(w io.Writer, opts *slog.HandlerOptions) *ColorTextHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}
	return &ColorTextHandler{w: w, opts: opts}
}

func (h *ColorTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColorTextHandler) Handle(_ context.Context, r slog.Record) error {
	levelStr, color := h.levelColor(r.Level)

	// 格式化时间
	timeStr := r.Time.Format("2006/01/02 15:04:05")

	// 日志头：[时间][等级]
	prefix := color + "[" + timeStr + "][" + levelStr + "]" + colorReset

	// 消息
	msg := r.Message

	// 合并输出
	line := prefix + " " + msg

	// 添加 attrs（key=value 形式）
	r.Attrs(func(attr slog.Attr) bool {
		line += " " + attr.Key + "=" + attr.Value.String()
		return true
	})
	line += "\n"

	_, err := h.w.Write([]byte(line))
	return err
}

func (h *ColorTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *ColorTextHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *ColorTextHandler) levelColor(level slog.Level) (string, string) {
	switch {
	case level < slog.LevelInfo:
		return "DEBUG", colorDebug
	case level < slog.LevelWarn:
		return "INFO ", colorInfo
	case level < slog.LevelError:
		return "WARN ", colorWarn
	default:
		return "ERROR", colorError
	}
}

func initLogger() {
	handler := NewColorTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

func GetLogger() *slog.Logger {
	loggerOnce.Do(initLogger)
	return defaultLogger
}
