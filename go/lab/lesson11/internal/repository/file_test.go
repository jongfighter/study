package repository

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"golab/lesson11/internal/domain"
	"golab/lesson11/internal/service"
)

var (
	_ service.Repository = (*Memory)(nil)
	_ service.Repository = (*File)(nil)
)

func TestFilePersistsAcrossRestart(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "tasks.json")

	first, err := NewFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := first.Create(ctx, "남아야 하는 할 일"); err != nil {
		t.Fatal(err)
	}
	removed, _ := first.Create(ctx, "지울 할 일")
	done := true
	if _, err := first.Update(ctx, 1, domain.UpdateInput{Done: &done}); err != nil {
		t.Fatal(err)
	}
	if err := first.Delete(ctx, removed.ID); err != nil {
		t.Fatal(err)
	}

	second, err := NewFile(path)
	if err != nil {
		t.Fatal(err)
	}
	page, _ := second.List(ctx, domain.ListFilter{})
	if len(page.Tasks) != 1 || page.Tasks[0].Title != "남아야 하는 할 일" || !page.Tasks[0].Done {
		t.Fatalf("다시 열었을 때 상태가 같아야 합니다: %+v", page.Tasks)
	}

	next, _ := second.Create(ctx, "다음 할 일")
	if next.ID != 3 {
		t.Errorf("지운 ID를 다시 쓰지 않고 3번이어야 합니다: %d", next.ID)
	}
}

func TestFileMissingFileStartsEmpty(t *testing.T) {
	f, err := NewFile(filepath.Join(t.TempDir(), "없는 파일.json"))
	if err != nil {
		t.Fatal(err)
	}
	page, _ := f.List(context.Background(), domain.ListFilter{})
	if len(page.Tasks) != 0 {
		t.Errorf("빈 상태로 시작해야 합니다: %+v", page.Tasks)
	}
}

func TestFileCorruptDataIsAnError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	if err := os.WriteFile(path, []byte("이건 JSON이 아닙니다"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := NewFile(path); err == nil {
		t.Error("깨진 파일은 조용히 무시하지 말고 에러를 돌려줘야 합니다")
	}
}

func TestFileLeavesNoTempFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	f, _ := NewFile(path)
	if _, err := f.Create(context.Background(), "임시 파일 확인"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("임시 파일이 남아 있습니다: %v", err)
	}
}

func TestFileFailsFastWhenDirectoryIsMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "없는폴더", "tasks.json")
	if _, err := NewFile(path); err == nil {
		t.Error("폴더가 없으면 시작할 때 에러가 나야 합니다")
	}
}
