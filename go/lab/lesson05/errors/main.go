package main

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("할 일을 찾을 수 없습니다")

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Reason
}

type Store struct {
	tasks map[int]string
}

func NewStore() *Store {
	return &Store{tasks: make(map[int]string)}
}

func (s *Store) Add(id int, title string) error {
	if strings.TrimSpace(title) == "" {
		return &ValidationError{Field: "title", Reason: "비어 있을 수 없습니다"}
	}
	s.tasks[id] = title
	return nil
}

func (s *Store) Get(id int) (string, error) {
	title, ok := s.tasks[id]
	if !ok {
		return "", fmt.Errorf("할 일 %d 조회: %w", id, ErrNotFound)
	}
	return title, nil
}

func main() {
	store := NewStore()

	if err := store.Add(1, "에러 처리 배우기"); err != nil {
		fmt.Println("추가 실패:", err)
	}

	title, err := store.Get(1)
	fmt.Println("조회:", title, err)

	_, err = store.Get(99)
	fmt.Println("오류:", err)
	fmt.Println("== 로 비교:", err == ErrNotFound)
	fmt.Println("errors.Is:", errors.Is(err, ErrNotFound))
	fmt.Println("Unwrap:", errors.Unwrap(err))

	err = store.Add(2, "   ")
	wrapped := fmt.Errorf("할 일 추가: %w", err)
	fmt.Println("오류:", wrapped)

	var validationErr *ValidationError
	if errors.As(wrapped, &validationErr) {
		fmt.Println("검증 실패 필드:", validationErr.Field)
	}

	switch {
	case errors.Is(wrapped, ErrNotFound):
		fmt.Println("404로 응답")
	case errors.As(wrapped, &validationErr):
		fmt.Println("400으로 응답")
	default:
		fmt.Println("500으로 응답")
	}
}
