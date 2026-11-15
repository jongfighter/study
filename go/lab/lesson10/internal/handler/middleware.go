package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

// Middleware는 핸들러를 감싸 공통 기능을 덧붙이는 함수입니다.
type Middleware func(http.Handler) http.Handler

// Chain은 h를 mws로 감쌉니다. 목록의 첫 번째가 가장 바깥쪽입니다.
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// requestIDKey는 context 키로만 쓰는 비공개 타입입니다. 다른 패키지의 키와 겹치지 않습니다.
type requestIDKey struct{}

// RequestIDFrom은 RequestID 미들웨어가 context에 넣어 둔 요청 ID를 꺼냅니다.
func RequestIDFrom(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey{}).(string)
	return id
}

// RequestID는 요청마다 무작위 ID를 만들어 context와 X-Request-Id 응답 헤더에 넣습니다.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := make([]byte, 8)
		_, _ = rand.Read(buf)
		id := hex.EncodeToString(buf)

		w.Header().Set("X-Request-Id", id)
		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Unwrap은 http.ResponseController가 원래 ResponseWriter에 접근할 수 있게 합니다.
func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

// Logging은 요청마다 구조화된 로그 한 줄을 남깁니다. 5xx는 ERROR, 4xx는 WARN, 나머지는 INFO입니다.
func Logging(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)

			level := slog.LevelInfo
			switch {
			case rec.status >= 500:
				level = slog.LevelError
			case rec.status >= 400:
				level = slog.LevelWarn
			}
			logger.Log(r.Context(), level, "request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration", time.Since(start),
				"request_id", RequestIDFrom(r.Context()),
			)
		})
	}
}

// Recover는 핸들러에서 난 패닉을 잡아 500 응답으로 바꿉니다.
func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					logger.ErrorContext(r.Context(), "패닉 복구", "panic", v, "request_id", RequestIDFrom(r.Context()))
					writeError(w, http.StatusInternalServerError, "서버 내부 오류")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
