package repository

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"golab/lesson11/internal/domain"
)

func TestMemoryCreateAndList(t *testing.T) {
	ctx := context.Background()
	fixed := time.Date(2026, 11, 15, 9, 0, 0, 0, time.UTC)
	m := NewMemory()
	m.now = func() time.Time { return fixed }

	first, _ := m.Create(ctx, "첫 번째")
	second, _ := m.Create(ctx, "두 번째")

	if first.ID != 1 || second.ID != 2 {
		t.Errorf("ID가 1, 2로 이어져야 합니다: %d, %d", first.ID, second.ID)
	}
	if !first.CreatedAt.Equal(fixed) {
		t.Errorf("CreatedAt = %v, want %v", first.CreatedAt, fixed)
	}

	page, err := m.List(ctx, domain.ListFilter{})
	if err != nil {
		t.Fatal(err)
	}
	tasks := page.Tasks
	if len(tasks) != 2 || tasks[0].Title != "첫 번째" || tasks[1].Title != "두 번째" {
		t.Errorf("목록이 ID 순서가 아닙니다: %+v", tasks)
	}
}

func TestMemoryNotFound(t *testing.T) {
	ctx := context.Background()
	m := NewMemory()
	done := true

	tests := map[string]func() error{
		"Get": func() error {
			_, err := m.Get(ctx, 99)
			return err
		},
		"Update": func() error {
			_, err := m.Update(ctx, 99, domain.UpdateInput{Done: &done})
			return err
		},
		"Delete": func() error {
			return m.Delete(ctx, 99)
		},
	}
	for name, call := range tests {
		t.Run(name, func(t *testing.T) {
			if err := call(); !errors.Is(err, domain.ErrNotFound) {
				t.Errorf("got %v, want ErrNotFound", err)
			}
		})
	}
}

func TestMemoryUpdateKeepsOtherFields(t *testing.T) {
	ctx := context.Background()
	m := NewMemory()
	created, _ := m.Create(ctx, "원래 제목")

	done := true
	updated, err := m.Update(ctx, created.ID, domain.UpdateInput{Done: &done})
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Done || updated.Title != "원래 제목" {
		t.Errorf("done만 바뀌어야 합니다: %+v", updated)
	}
}

func TestMemoryConcurrentCreate(t *testing.T) {
	ctx := context.Background()
	m := NewMemory()

	const n = 50
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := m.Create(ctx, "동시 생성"); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	page, _ := m.List(ctx, domain.ListFilter{})
	tasks := page.Tasks
	if len(tasks) != n {
		t.Fatalf("len = %d, want %d", len(tasks), n)
	}
	seen := make(map[int]bool)
	for _, task := range tasks {
		if seen[task.ID] {
			t.Fatalf("중복된 ID: %d", task.ID)
		}
		seen[task.ID] = true
	}
}
