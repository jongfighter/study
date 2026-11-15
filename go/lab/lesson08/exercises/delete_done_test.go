package exercises

import (
	"errors"
	"sync"
	"testing"
)

func TestRemoveOnlyDone(t *testing.T) {
	s := NewStore()
	s.Add(1, false)
	s.Add(2, true)

	if err := s.Remove(1); !errors.Is(err, ErrNotDone) {
		t.Errorf("완료 전에는 ErrNotDone이어야 합니다: %v", err)
	}
	if err := s.Remove(2); err != nil {
		t.Errorf("완료된 할 일은 지워져야 합니다: %v", err)
	}
	if err := s.Remove(2); !errors.Is(err, ErrNotFound) {
		t.Errorf("이미 지운 id는 ErrNotFound여야 합니다: %v", err)
	}
	if err := s.Remove(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("없는 id는 ErrNotFound여야 합니다: %v", err)
	}
}

func TestRemoveConcurrentlyDeletesOnce(t *testing.T) {
	s := NewStore()
	s.Add(1, true)

	const n = 50
	results := make(chan error, n)
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- s.Remove(1)
		}()
	}
	wg.Wait()
	close(results)

	succeeded := 0
	for err := range results {
		if err == nil {
			succeeded++
		} else if !errors.Is(err, ErrNotFound) {
			t.Errorf("예상하지 못한 에러: %v", err)
		}
	}
	if succeeded != 1 {
		t.Errorf("삭제는 정확히 한 번만 성공해야 합니다: %d번 성공", succeeded)
	}
}
