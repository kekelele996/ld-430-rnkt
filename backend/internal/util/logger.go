package util

import (
	"log/slog"
	"os"
)

// NewLogger builds a structured slog.Logger with JSON output.
func NewLogger(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
