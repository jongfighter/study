package handler

import (
	"context"
	"net/http"
	"testing"

	"golab/lesson11/internal/domain"
)

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
			h := New(fakeTasks{update: update}, testLogger).Routes()

			rec := do(t, h, http.MethodPatch, tt.path, tt.body)

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
			h := New(fakeTasks{remove: remove}, testLogger).Routes()

			rec := do(t, h, http.MethodDelete, tt.path, "")

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantStatus == http.StatusNoContent && rec.Body.Len() != 0 {
				t.Errorf("204 응답에는 본문이 없어야 합니다: %q", rec.Body)
			}
		})
	}
}
