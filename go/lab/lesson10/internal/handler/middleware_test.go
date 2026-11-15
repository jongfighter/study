package handler

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newBufferLogger는 시간 필드를 뺀 텍스트 로그를 buf에 쌓는 로거를 만듭니다.
func newBufferLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
}

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

func TestRequestIDIsSetAndStored(t *testing.T) {
	var fromContext string
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fromContext = RequestIDFrom(r.Context())
	})

	rec := httptest.NewRecorder()
	Chain(inner, RequestID).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	header := rec.Header().Get("X-Request-Id")
	if len(header) != 16 {
		t.Errorf("X-Request-Id = %q, 16자리 16진수여야 합니다", header)
	}
	if fromContext != header {
		t.Errorf("context의 ID %q와 헤더의 ID %q가 같아야 합니다", fromContext, header)
	}
}

func TestLoggingLevelByStatus(t *testing.T) {
	tests := []struct {
		status     int
		wantStatus string
		wantLevel  string
	}{
		{http.StatusOK, "status=200", "level=INFO"},
		{http.StatusNotFound, "status=404", "level=WARN"},
		{http.StatusInternalServerError, "status=500", "level=ERROR"},
	}
	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			var buf bytes.Buffer
			inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.status)
			})

			Chain(inner, RequestID, Logging(newBufferLogger(&buf))).
				ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/tasks", nil))

			out := buf.String()
			for _, want := range []string{tt.wantLevel, tt.wantStatus, "msg=request", "method=GET", "path=/tasks", "request_id="} {
				if !strings.Contains(out, want) {
					t.Errorf("로그 %q에 %q가 없습니다", out, want)
				}
			}
		})
	}
}

func TestRecoverTurnsPanicInto500(t *testing.T) {
	var buf bytes.Buffer
	boom := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	rec := httptest.NewRecorder()
	Chain(boom, Recover(newBufferLogger(&buf))).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "서버 내부 오류") {
		t.Errorf("body = %s", rec.Body)
	}
	if out := buf.String(); !strings.Contains(out, "level=ERROR") || !strings.Contains(out, "panic=boom") {
		t.Errorf("에러 로그에 패닉 내용이 남아야 합니다: %q", out)
	}
}
