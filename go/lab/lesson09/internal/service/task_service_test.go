package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golab/lesson09/internal/domain"
	"golab/lesson09/internal/repository"
)

func newTestService() *Service {
	return New(repository.NewMemory())
}

func TestCreateValidation(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantTitle string
		wantErr   bool
	}{
		{"정상", "Go 공부하기", "Go 공부하기", false},
		{"앞뒤 공백 제거", "  장보기  ", "장보기", false},
		{"빈 문자열", "", "", true},
		{"공백만", "   ", "", true},
		{"한글 100자는 허용", strings.Repeat("가", 100), strings.Repeat("가", 100), false},
		{"한글 101자는 거부", strings.Repeat("가", 101), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTestService().Create(context.Background(), tt.title)

			var validation *domain.ValidationError
			if tt.wantErr {
				if !errors.As(err, &validation) {
					t.Fatalf("ValidationError가 와야 합니다: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("예상하지 못한 에러: %v", err)
			}
			if got.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", got.Title, tt.wantTitle)
			}
		})
	}
}

func TestUpdateRules(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()
	created, _ := svc.Create(ctx, "원래 제목")

	blank := "  "
	done := true

	var validation *domain.ValidationError
	if _, err := svc.Update(ctx, created.ID, domain.UpdateInput{}); !errors.As(err, &validation) {
		t.Errorf("수정할 필드가 없으면 ValidationError여야 합니다: %v", err)
	}
	if _, err := svc.Update(ctx, created.ID, domain.UpdateInput{Title: &blank}); !errors.As(err, &validation) {
		t.Errorf("빈 제목은 ValidationError여야 합니다: %v", err)
	}

	got, err := svc.Update(ctx, created.ID, domain.UpdateInput{Done: &done})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Done {
		t.Error("Done이 true로 바뀌어야 합니다")
	}
}

func TestGetNotFound(t *testing.T) {
	_, err := newTestService().Get(context.Background(), 42)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
	}
}

type failingRepo struct {
	Repository
}

func (failingRepo) Create(context.Context, string) (domain.Task, error) {
	return domain.Task{}, errors.New("db down")
}

func TestRepositoryErrorsPropagate(t *testing.T) {
	svc := New(failingRepo{})
	_, err := svc.Create(context.Background(), "아무 제목")
	if err == nil || err.Error() != "db down" {
		t.Errorf("저장소 에러가 그대로 전달되어야 합니다: %v", err)
	}
}
