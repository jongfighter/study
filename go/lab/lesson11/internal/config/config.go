// Package config는 환경변수에서 서버 설정을 읽습니다.
package config

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

// Config는 서버 실행에 필요한 설정입니다.
type Config struct {
	Port            string
	LogLevel        slog.Level
	LogFormat       string
	ShutdownTimeout time.Duration
	DataFile        string // 비어 있으면 메모리에만 저장합니다.
}

// Load는 getenv로 환경변수를 읽어 Config를 만듭니다.
// 값이 없으면 기본값을 쓰고, 값이 잘못되었으면 에러를 돌려줍니다.
// getenv를 인자로 받으므로 테스트에서 os.Getenv 대신 가짜 함수를 넘길 수 있습니다.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:            "8080",
		LogLevel:        slog.LevelInfo,
		LogFormat:       "text",
		ShutdownTimeout: 10 * time.Second,
	}

	if v := getenv("PORT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 || n > 65535 {
			return Config{}, fmt.Errorf("PORT는 1~65535 사이의 숫자여야 합니다: %q", v)
		}
		cfg.Port = v
	}

	if v := getenv("LOG_LEVEL"); v != "" {
		if err := cfg.LogLevel.UnmarshalText([]byte(v)); err != nil {
			return Config{}, fmt.Errorf("LOG_LEVEL이 올바르지 않습니다: %q (debug, info, warn, error 중 하나)", v)
		}
	}

	if v := getenv("LOG_FORMAT"); v != "" {
		format := strings.ToLower(v)
		if format != "text" && format != "json" {
			return Config{}, fmt.Errorf("LOG_FORMAT은 text 또는 json이어야 합니다: %q", v)
		}
		cfg.LogFormat = format
	}

	if v := getenv("SHUTDOWN_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT은 0보다 큰 시간이어야 합니다(예: 15s): %q", v)
		}
		cfg.ShutdownTimeout = d
	}

	cfg.DataFile = getenv("DATA_FILE")

	return cfg, nil
}
