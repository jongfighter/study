package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golab/lesson09/internal/domain"
	"golab/lesson09/internal/repository"
	"golab/lesson09/internal/service"
)

// fakeTasks는 테스트용 가짜 TaskService입니다. 필요한 동작만 함수로 채웁니다.
type fakeTasks struct {
	create func(ctx context.Context, title string) (domain.Task, error)
	get    func(ctx context.Context, id int) (domain.Task, error)
	update func(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error)
	remove func(ctx context.Context, id int) error
}

var errUnused = errors.New("이 테스트에서는 쓰지 않는 메서드입니다")

func (f fakeTasks) List(context.Context) ([]domain.Task, error) {
	return nil, errUnused
}

func (f fakeTasks) Create(ctx context.Context, title string) (domain.Task, error) {
	return f.create(ctx, title)
}

func (f fakeTasks) Get(ctx context.Context, id int) (domain.Task, error) {
	return f.get(ctx, id)
}

func (f fakeTasks) Update(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error) {
	return f.update(ctx, id, in)
}

func (f fakeTasks) Delete(ctx context.Context, id int) error {
	return f.remove(ctx, id)
}

func do(t *testing.T, h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestCreateTask(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		create     func(ctx context.Context, title string) (domain.Task, error)
		wantStatus int
		wantBody   string
	}{
		{
			name: "성공",
			body: `{"title":"Go 공부하기"}`,
			create: func(_ context.Context, title string) (domain.Task, error) {
				return domain.Task{ID: 7, Title: title}, nil
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"id":7`,
		},
		{
			name: "검증 실패는 400",
			body: `{"title":""}`,
			create: func(context.Context, string) (domain.Task, error) {
				return domain.Task{}, &domain.ValidationError{Message: "title은 비어 있을 수 없습니다"}
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"title은 비어 있을 수 없습니다"}`,
		},
		{
			name: "깨진 JSON은 서비스를 부르지 않고 400",
			body: `{`,
			create: func(context.Context, string) (domain.Task, error) {
				t.Error("서비스가 호출되면 안 됩니다")
				return domain.Task{}, nil
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "요청 본문을 읽을 수 없습니다",
		},
		{
			name: "알 수 없는 필드는 400",
			body: `{"titel":"오타"}`,
			create: func(context.Context, string) (domain.Task, error) {
				t.Error("서비스가 호출되면 안 됩니다")
				return domain.Task{}, nil
			},
			wantStatus: http.StatusBadRequest,
			wantBody:   "unknown field",
		},
		{
			name: "예상하지 못한 에러는 500이고 내용을 숨김",
			body: `{"title":"x"}`,
			create: func(context.Context, string) (domain.Task, error) {
				return domain.Task{}, errors.New("db down")
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"서버 내부 오류"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(fakeTasks{create: tt.create}).Routes()

			rec := do(t, h, http.MethodPost, "/tasks", tt.body)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body = %s, want %q 포함", rec.Body, tt.wantBody)
			}
			if tt.wantStatus == http.StatusCreated {
				if got := rec.Header().Get("Location"); got != "/tasks/7" {
					t.Errorf("Location = %q, want /tasks/7", got)
				}
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	get := func(_ context.Context, id int) (domain.Task, error) {
		if id == 1 {
			return domain.Task{ID: 1, Title: "있는 할 일"}, nil
		}
		return domain.Task{}, domain.ErrNotFound
	}
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"성공", "/tasks/1", http.StatusOK},
		{"없는 id", "/tasks/2", http.StatusNotFound},
		{"숫자가 아닌 id", "/tasks/abc", http.StatusBadRequest},
		{"0 이하의 id", "/tasks/0", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(fakeTasks{get: get}).Routes()

			rec := do(t, h, http.MethodGet, tt.path, "")

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}

// TestTaskFlow는 가짜 없이 service와 repository까지 모두 연결해 흐름 전체를 확인합니다.
func TestTaskFlow(t *testing.T) {
	h := New(service.New(repository.NewMemory()))
	srv := httptest.NewServer(Chain(h.Routes(), Recover))
	defer srv.Close()

	send := func(method, path, body string) *http.Response {
		t.Helper()
		req, err := http.NewRequest(method, srv.URL+path, strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { resp.Body.Close() })
		return resp
	}

	resp := send(http.MethodPost, "/tasks", `{"title":"흐름 확인"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST status = %d", resp.StatusCode)
	}
	var created domain.Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	if resp := send(http.MethodGet, "/tasks/1", ""); resp.StatusCode != http.StatusOK {
		t.Errorf("GET status = %d", resp.StatusCode)
	}

	resp = send(http.MethodPatch, "/tasks/1", `{"done":true}`)
	var patched domain.Task
	if err := json.NewDecoder(resp.Body).Decode(&patched); err != nil {
		t.Fatal(err)
	}
	if !patched.Done || patched.Title != created.Title {
		t.Errorf("PATCH 결과가 이상합니다: %+v", patched)
	}

	if resp := send(http.MethodDelete, "/tasks/1", ""); resp.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE status = %d", resp.StatusCode)
	}
	if resp := send(http.MethodGet, "/tasks/1", ""); resp.StatusCode != http.StatusNotFound {
		t.Errorf("삭제 뒤 GET status = %d, want 404", resp.StatusCode)
	}

	resp = send(http.MethodPut, "/tasks/1", "")
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("PUT status = %d, want 405", resp.StatusCode)
	}
	if allow := resp.Header.Get("Allow"); !strings.Contains(allow, "PATCH") {
		t.Errorf("Allow 헤더에 PATCH가 있어야 합니다: %q", allow)
	}
}
