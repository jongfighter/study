package exercises

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMaxBodyRejectsByContentLength(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(strings.Repeat("x", 100)))
	rec := httptest.NewRecorder()

	MaxBody(10)(next).ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("status = %d, want 413", rec.Code)
	}
	if called {
		t.Error("제한을 넘으면 다음 핸들러를 부르면 안 됩니다")
	}
}

func TestMaxBodyCutsUnknownLengthBody(t *testing.T) {
	var readErr error
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, readErr = io.ReadAll(r.Body)
	})
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(strings.Repeat("x", 100)))
	req.ContentLength = -1

	MaxBody(10)(next).ServeHTTP(httptest.NewRecorder(), req)

	if !IsTooLarge(readErr) {
		t.Errorf("읽는 도중 MaxBytesError가 나와야 합니다: %v", readErr)
	}
}

func TestMaxBodyAllowsSmallBody(t *testing.T) {
	var got string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		got = string(b)
	})
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("작은 본문"))

	MaxBody(1024)(next).ServeHTTP(httptest.NewRecorder(), req)

	if got != "작은 본문" {
		t.Errorf("got %q", got)
	}
}
