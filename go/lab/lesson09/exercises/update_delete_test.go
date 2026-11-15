package exercises

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golab/lesson09/internal/domain"
	"golab/lesson09/internal/handler"
)

// stubService는 PATCH와 DELETE만 흉내 내는 가짜 TaskService입니다.
type stubService struct {
	update func(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error)
	remove func(ctx context.Context, id int) error
}

var errUnused = errors.New("이 테스트에서는 쓰지 않는 메서드입니다")

func (s stubService) List(context.Context) ([]domain.Task, error) {
	return nil, errUnused
}

func (s stubService) Create(context.Context, string) (domain.Task, error) {
	return domain.Task{}, errUnused
}

func (s stubService) Get(context.Context, int) (domain.Task, error) {
	return domain.Task{}, errUnused
}

func (s stubService) Update(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error) {
	return s.update(ctx, id, in)
}

func (s stubService) Delete(ctx context.Context, id int) error {
	return s.remove(ctx, id)
}

func send(h http.Handler, method, target, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestUpdateTask(t *testing.T) {
	update := func(_ context.Context, id int, in domain.UpdateInput) (domain.Task, error) {
		switch {
		case id == 99:
			return domain.Task{}, domain.ErrNotFound
		case in.Title == nil && in.Done == nil:
			return domain.Task{}, &domain.ValidationError{Message: "수정할 필드(title, done)가 없습니다"}
		}
		return domain.Task{ID: id, Title: "수정됨", Done: in.Done != nil && *in.Done}, nil
	}
	tests := []struct {
		name       string
		path       string
		body       string
		wantStatus int
	}{
		{"done 수정", "/tasks/1", `{"done":true}`, http.StatusOK},
		{"수정할 필드 없음", "/tasks/1", `{}`, http.StatusBadRequest},
		{"없는 id", "/tasks/99", `{"done":true}`, http.StatusNotFound},
		{"숫자가 아닌 id", "/tasks/x", `{"done":true}`, http.StatusBadRequest},
		{"깨진 JSON", "/tasks/1", `{"done":`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.New(stubService{update: update}).Routes()

			rec := send(h, http.MethodPatch, tt.path, tt.body)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	remove := func(_ context.Context, id int) error {
		if id == 1 {
			return nil
		}
		return domain.ErrNotFound
	}
	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{"성공은 204", "/tasks/1", http.StatusNoContent},
		{"없는 id는 404", "/tasks/2", http.StatusNotFound},
		{"숫자가 아닌 id는 400", "/tasks/abc", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.New(stubService{remove: remove}).Routes()

			rec := send(h, http.MethodDelete, tt.path, "")

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusNoContent && rec.Body.Len() != 0 {
				t.Errorf("204 응답에는 본문이 없어야 합니다: %q", rec.Body)
			}
		})
	}
}
