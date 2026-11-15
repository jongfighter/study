package config

import (
	"log/slog"
	"testing"
	"time"
)

func fakeEnv(values map[string]string) func(string) string {
	return func(key string) string { return values[key] }
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load(fakeEnv(nil))
	if err != nil {
		t.Fatal(err)
	}
	want := Config{
		Port:            "8080",
		LogLevel:        slog.LevelInfo,
		LogFormat:       "text",
		ShutdownTimeout: 10 * time.Second,
	}
	if cfg != want {
		t.Errorf("got %+v, want %+v", cfg, want)
	}
}

func TestLoadOverrides(t *testing.T) {
	cfg, err := Load(fakeEnv(map[string]string{
		"PORT":             "9090",
		"LOG_LEVEL":        "debug",
		"LOG_FORMAT":       "JSON",
		"SHUTDOWN_TIMEOUT": "3s",
		"DATA_FILE":        "tasks.json",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "9090" || cfg.LogLevel != slog.LevelDebug || cfg.LogFormat != "json" || cfg.ShutdownTimeout != 3*time.Second || cfg.DataFile != "tasks.json" {
		t.Errorf("설정이 반영되지 않았습니다: %+v", cfg)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
	}{
		{"PORT가 숫자가 아님", map[string]string{"PORT": "abc"}},
		{"PORT가 범위 밖", map[string]string{"PORT": "70000"}},
		{"LOG_LEVEL이 알 수 없는 값", map[string]string{"LOG_LEVEL": "loud"}},
		{"LOG_FORMAT이 알 수 없는 값", map[string]string{"LOG_FORMAT": "xml"}},
		{"SHUTDOWN_TIMEOUT 형식 오류", map[string]string{"SHUTDOWN_TIMEOUT": "열 초"}},
		{"SHUTDOWN_TIMEOUT이 0 이하", map[string]string{"SHUTDOWN_TIMEOUT": "0s"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Load(fakeEnv(tt.env)); err == nil {
				t.Error("에러가 와야 합니다")
			}
		})
	}
}
