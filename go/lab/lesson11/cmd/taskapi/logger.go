package main

import (
	"io"
	"log/slog"

	"golab/lesson11/internal/config"
)

// newLogger는 설정에 맞는 slog.Logger를 만듭니다.
func newLogger(w io.Writer, cfg config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: cfg.LogLevel}
	if cfg.LogFormat == "json" {
		return slog.New(slog.NewJSONHandler(w, opts))
	}
	return slog.New(slog.NewTextHandler(w, opts))
}
