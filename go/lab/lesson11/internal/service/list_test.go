package service

import (
	"context"
	"errors"
	"testing"

	"golab/lesson11/internal/domain"
)

func TestListDefaultsAndValidation(t *testing.T) {
	ctx := context.Background()
	svc := newTestService()
	for range 25 {
		if _, err := svc.Create(ctx, "Go 공부"); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("limit을 생략하면 기본값 20", func(t *testing.T) {
		page, err := svc.List(ctx, domain.ListFilter{})
		if err != nil {
			t.Fatal(err)
		}
		if len(page.Tasks) != 20 || page.Total != 25 {
			t.Errorf("len = %d, total = %d", len(page.Tasks), page.Total)
		}
	})

	t.Run("검색어 앞뒤 공백은 무시", func(t *testing.T) {
		page, err := svc.List(ctx, domain.ListFilter{Query: "  go  ", Limit: 5})
		if err != nil {
			t.Fatal(err)
		}
		if page.Total != 25 {
			t.Errorf("Total = %d, want 25", page.Total)
		}
	})

	tests := []struct {
		name    string
		filter  domain.ListFilter
		wantErr bool
	}{
		{"limit이 음수", domain.ListFilter{Limit: -1}, true},
		{"limit이 100 초과", domain.ListFilter{Limit: 101}, true},
		{"offset이 음수", domain.ListFilter{Offset: -1}, true},
		{"limit이 상한과 같으면 통과", domain.ListFilter{Limit: 100}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.List(ctx, tt.filter)

			var validation *domain.ValidationError
			if tt.wantErr && !errors.As(err, &validation) {
				t.Errorf("ValidationError가 와야 합니다: %v", err)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("예상하지 못한 에러: %v", err)
			}
		})
	}
}
