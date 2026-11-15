// Package exercises는 레슨 09 연습 문제의 정답 예시입니다.
package exercises

import (
	"crypto/subtle"
	"net/http"
)

// RequireAPIKey는 X-API-Key 헤더가 key와 같은 요청만 통과시킵니다.
// GET /health는 상태 확인용이므로 키 없이도 통과시킵니다.
func RequireAPIKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/health" {
				next.ServeHTTP(w, r)
				return
			}
			got := r.Header.Get("X-API-Key")
			if subtle.ConstantTimeCompare([]byte(got), []byte(key)) != 1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"인증이 필요합니다"}` + "\n"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
