// Package repository는 Task를 저장하는 구현을 모아 둡니다.
package repository

import (
	"context"
	"sort"
	"sync"
	"time"

	"golab/lesson09/internal/domain"
)

// Memory는 Task를 메모리에 보관하는 저장소입니다. 여러 고루틴이 동시에 써도 안전합니다.
type Memory struct {
	mu     sync.Mutex
	nextID int
	tasks  map[int]domain.Task
	now    func() time.Time
}

// NewMemory는 빈 저장소를 만듭니다.
func NewMemory() *Memory {
	return &Memory{nextID: 1, tasks: make(map[int]domain.Task), now: time.Now}
}

func (m *Memory) List(_ context.Context) ([]domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (m *Memory) Create(_ context.Context, title string) (domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t := domain.Task{ID: m.nextID, Title: title, CreatedAt: m.now().UTC()}
	m.nextID++
	m.tasks[t.ID] = t
	return t, nil
}

func (m *Memory) Get(_ context.Context, id int) (domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, nil
}

func (m *Memory) Update(_ context.Context, id int, in domain.UpdateInput) (domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	if in.Title != nil {
		t.Title = *in.Title
	}
	if in.Done != nil {
		t.Done = *in.Done
	}
	m.tasks[id] = t
	return t, nil
}

func (m *Memory) Delete(_ context.Context, id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}
