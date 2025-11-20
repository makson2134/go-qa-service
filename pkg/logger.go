package pkg

import (
	"log/slog"
	"os"
	"strings"
)

// Warning level and json format by default
func NewLogger(level, format string) *slog.Logger {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelWarn
	}

	var handler slog.Handler

	switch strings.ToLower(format) {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	return slog.New(handler)
}
