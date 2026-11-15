package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"golab/lesson11/internal/domain"
)

// File은 Memory 위에 JSON 파일 저장을 얹은 저장소입니다.
// 서버를 다시 시작해도 Task가 남고, 변경할 때마다 파일 전체를 다시 씁니다.
type File struct {
	mu   sync.Mutex
	path string
	mem  *Memory
}

type fileData struct {
	NextID int           `json:"next_id"`
	Tasks  []domain.Task `json:"tasks"`
}

// NewFile은 path의 데이터를 읽어 저장소를 만듭니다. 파일이 없으면 빈 상태로 시작합니다.
// 시작할 때 한 번 저장해 보므로, 폴더가 없거나 쓸 수 없으면 서버가 뜨기 전에 에러가 납니다.
func NewFile(path string) (*File, error) {
	f := &File{path: path, mem: NewMemory()}
	if err := f.load(); err != nil {
		return nil, err
	}
	if err := f.save(); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *File) load() error {
	b, err := os.ReadFile(f.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("데이터 파일을 읽을 수 없습니다: %w", err)
	}

	var data fileData
	if err := json.Unmarshal(b, &data); err != nil {
		return fmt.Errorf("데이터 파일 형식이 올바르지 않습니다: %w", err)
	}
	f.mem.restore(data.NextID, data.Tasks)
	return nil
}

// save는 임시 파일에 쓴 뒤 이름을 바꿔, 쓰다가 멈춰도 기존 파일이 깨지지 않게 합니다.
func (f *File) save() error {
	nextID, tasks := f.mem.snapshot()
	b, err := json.MarshalIndent(fileData{NextID: nextID, Tasks: tasks}, "", "  ")
	if err != nil {
		return err
	}

	tmp := f.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return fmt.Errorf("데이터 파일을 쓸 수 없습니다: %w", err)
	}
	if err := os.Rename(tmp, f.path); err != nil {
		return fmt.Errorf("데이터 파일을 교체할 수 없습니다: %w", err)
	}
	return nil
}

func (f *File) List(ctx context.Context, filter domain.ListFilter) (domain.Page, error) {
	return f.mem.List(ctx, filter)
}

func (f *File) Get(ctx context.Context, id int) (domain.Task, error) {
	return f.mem.Get(ctx, id)
}

func (f *File) Create(ctx context.Context, title string) (domain.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	t, err := f.mem.Create(ctx, title)
	if err != nil {
		return domain.Task{}, err
	}
	return t, f.save()
}

func (f *File) Update(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	t, err := f.mem.Update(ctx, id, in)
	if err != nil {
		return domain.Task{}, err
	}
	return t, f.save()
}

func (f *File) Delete(ctx context.Context, id int) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.mem.Delete(ctx, id); err != nil {
		return err
	}
	return f.save()
}
