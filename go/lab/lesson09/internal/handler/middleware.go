package handler

import (
	"log"
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

// Logging은 요청마다 메서드, 경로, 상태 코드, 걸린 시간을 기록합니다.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}

// Recover는 핸들러에서 난 패닉을 잡아 500 응답으로 바꿉니다.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				log.Printf("패닉 복구: %v", v)
				writeError(w, http.StatusInternalServerError, "서버 내부 오류")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
