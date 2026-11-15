// Package repository는 Task를 저장하는 구현을 모아 둡니다.
package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"golab/lesson11/internal/domain"
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

func (m *Memory) List(_ context.Context, f domain.ListFilter) (domain.Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	query := strings.ToLower(f.Query)
	matched := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		if f.Done != nil && t.Done != *f.Done {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(t.Title), query) {
			continue
		}
		matched = append(matched, t)
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].ID < matched[j].ID })

	return domain.Page{Tasks: paginate(matched, f.Offset, f.Limit), Total: len(matched)}, nil
}

// paginate는 tasks에서 offset만큼 건너뛰고 최대 limit개를 돌려줍니다. limit이 0 이하이면 끝까지 돌려줍니다.
func paginate(tasks []domain.Task, offset, limit int) []domain.Task {
	if offset >= len(tasks) {
		return []domain.Task{}
	}
	tasks = tasks[offset:]
	if limit > 0 && limit < len(tasks) {
		tasks = tasks[:limit]
	}
	return tasks
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

// snapshot은 다음 ID와 모든 Task를 ID 순서로 복사해 돌려줍니다.
func (m *Memory) snapshot() (int, []domain.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	tasks := make([]domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		tasks = append(tasks, t)
	}
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].ID < tasks[j].ID })
	return m.nextID, tasks
}

// restore는 저장된 상태로 내용을 통째로 바꿉니다.
func (m *Memory) restore(nextID int, tasks []domain.Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.nextID = nextID
	m.tasks = make(map[int]domain.Task, len(tasks))
	for _, t := range tasks {
		m.tasks[t.ID] = t
		if t.ID >= m.nextID {
			m.nextID = t.ID + 1
		}
	}
}
