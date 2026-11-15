package repository

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"golab/lesson11/internal/domain"
	"golab/lesson11/internal/service"
)

// runContract는 어떤 저장소 구현이든 지켜야 하는 약속을 확인합니다.
func runContract(t *testing.T, newRepo func(t *testing.T) service.Repository) {
	ctx := context.Background()

	t.Run("만든 Task를 다시 읽을 수 있다", func(t *testing.T) {
		repo := newRepo(t)
		created, err := repo.Create(ctx, "계약 테스트")
		if err != nil {
			t.Fatal(err)
		}
		got, err := repo.Get(ctx, created.ID)
		if err != nil {
			t.Fatal(err)
		}
		if got.Title != "계약 테스트" || got.Done {
			t.Errorf("got %+v", got)
		}
	})

	t.Run("없는 Task는 ErrNotFound", func(t *testing.T) {
		repo := newRepo(t)
		done := true
		if _, err := repo.Get(ctx, 99); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Get: %v", err)
		}
		if _, err := repo.Update(ctx, 99, domain.UpdateInput{Done: &done}); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Update: %v", err)
		}
		if err := repo.Delete(ctx, 99); !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Delete: %v", err)
		}
	})

	t.Run("목록은 ID 순서이고 필터와 쪽 나누기를 지킨다", func(t *testing.T) {
		repo := newRepo(t)
		for _, title := range []string{"가", "나", "다"} {
			if _, err := repo.Create(ctx, title); err != nil {
				t.Fatal(err)
			}
		}
		page, err := repo.List(ctx, domain.ListFilter{Limit: 2})
		if err != nil {
			t.Fatal(err)
		}
		if page.Total != 3 || len(page.Tasks) != 2 || page.Tasks[0].Title != "가" {
			t.Errorf("got %+v", page)
		}
	})
}

func TestMemoryContract(t *testing.T) {
	runContract(t, func(t *testing.T) service.Repository { return NewMemory() })
}

func TestFileContract(t *testing.T) {
	runContract(t, func(t *testing.T) service.Repository {
		repo, err := NewFile(filepath.Join(t.TempDir(), "tasks.json"))
		if err != nil {
			t.Fatal(err)
		}
		return repo
	})
}
