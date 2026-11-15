package handler

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChainOrder(t *testing.T) {
	var calls []string
	tag := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls = append(calls, name+" 시작")
				next.ServeHTTP(w, r)
				calls = append(calls, name+" 끝")
			})
		}
	}
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, "핸들러")
	})

	Chain(final, tag("A"), tag("B")).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	got := strings.Join(calls, ", ")
	want := "A 시작, B 시작, 핸들러, B 끝, A 끝"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRecoverTurnsPanicInto500(t *testing.T) {
	var logs bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(prev) })

	boom := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	rec := httptest.NewRecorder()
	Chain(boom, Recover).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "서버 내부 오류") {
		t.Errorf("body = %s", rec.Body)
	}
	if !strings.Contains(logs.String(), "패닉 복구: boom") {
		t.Errorf("로그에 패닉 내용이 남아야 합니다: %q", logs.String())
	}
}

func TestLoggingRecordsStatus(t *testing.T) {
	var logs bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(prev) })

	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	Chain(notFound, Logging).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/nope", nil))

	if !strings.Contains(logs.String(), "GET /nope -> 404") {
		t.Errorf("로그 = %q", logs.String())
	}
}
