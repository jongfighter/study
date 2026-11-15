// Package exercises는 레슨 10 연습 문제의 정답 예시입니다.
package exercises

import (
	"errors"
	"net/http"
)

// MaxBody는 요청 본문을 limit바이트까지만 읽도록 제한합니다.
// Content-Length가 이미 limit보다 크면 핸들러를 부르지 않고 413으로 응답하고,
// 길이를 모르는 요청은 http.MaxBytesReader로 감싸 읽는 도중에 끊습니다.
func MaxBody(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > limit {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				_, _ = w.Write([]byte(`{"error":"요청 본문이 너무 큽니다"}` + "\n"))
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

// IsTooLarge는 body를 읽다가 난 에러가 크기 제한 때문인지 알려 줍니다.
func IsTooLarge(err error) bool {
	var tooLarge *http.MaxBytesError
	return errors.As(err, &tooLarge)
}
