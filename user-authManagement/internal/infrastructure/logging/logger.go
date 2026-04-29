// Package logging configures the application's slog logger.
package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// New returns a configured slog.Logger.
// format is "json" or "text"; level is "debug"/"info"/"warn"/"error".
func New(level, format string) *slog.Logger {
	return NewWithWriter(os.Stdout, level, format)
}

func NewWithWriter(w io.Writer, level, format string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: lvl, AddSource: false}
	var h slog.Handler
	if strings.ToLower(format) == "text" {
		h = slog.NewTextHandler(w, opts)
	} else {
		h = slog.NewJSONHandler(w, opts)
	}
	return slog.New(h)
}
