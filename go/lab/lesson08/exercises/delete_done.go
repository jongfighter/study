package exercises

import (
	"errors"
	"sync"
)

var (
	// ErrNotFound는 해당 id의 할 일이 없을 때의 에러입니다. 핸들러에서 404로 바꿉니다.
	ErrNotFound = errors.New("할 일을 찾을 수 없습니다")
	// ErrNotDone은 완료되지 않은 할 일을 지우려 할 때의 에러입니다. 핸들러에서 409로 바꿉니다.
	ErrNotDone = errors.New("완료된 할 일만 삭제할 수 있습니다")
)

// Store는 id와 완료 여부만 가진 아주 작은 저장소입니다.
type Store struct {
	mu    sync.Mutex
	items map[int]bool
}

// NewStore는 빈 저장소를 만듭니다.
func NewStore() *Store {
	return &Store{items: make(map[int]bool)}
}

// Add는 id에 해당하는 할 일을 완료 여부와 함께 저장합니다.
func (s *Store) Add(id int, done bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = done
}

// Remove는 완료된 할 일만 지웁니다. 확인과 삭제를 같은 락 안에서 하므로
// 그 사이에 다른 요청이 끼어들 수 없습니다.
func (s *Store) Remove(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	done, ok := s.items[id]
	switch {
	case !ok:
		return ErrNotFound
	case !done:
		return ErrNotDone
	}
	delete(s.items, id)
	return nil
}
