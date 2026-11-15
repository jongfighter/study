package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"golab/lesson10/internal/config"
)

func TestNewLoggerText(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(&buf, config.Config{LogFormat: "text", LogLevel: slog.LevelInfo})

	logger.Info("안녕", "key", "value")

	out := buf.String()
	if !strings.Contains(out, "level=INFO") || !strings.Contains(out, "key=value") {
		t.Errorf("텍스트 형식이 아닙니다: %q", out)
	}
}

func TestNewLoggerJSON(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(&buf, config.Config{LogFormat: "json", LogLevel: slog.LevelInfo})

	logger.Info("안녕", "key", "value")

	var line map[string]any
	if err := json.Unmarshal(buf.Bytes(), &line); err != nil {
		t.Fatalf("JSON 한 줄이어야 합니다: %v (%q)", err, buf.String())
	}
	if line["level"] != "INFO" || line["msg"] != "안녕" || line["key"] != "value" {
		t.Errorf("필드가 이상합니다: %v", line)
	}
}

func TestNewLoggerRespectsLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := newLogger(&buf, config.Config{LogFormat: "text", LogLevel: slog.LevelWarn})

	logger.Info("보이지 않아야 함")
	logger.Warn("보여야 함")

	out := buf.String()
	if strings.Contains(out, "보이지 않아야 함") || !strings.Contains(out, "보여야 함") {
		t.Errorf("WARN 이상만 나와야 합니다: %q", out)
	}
}
