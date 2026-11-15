// Package domain은 Task API의 핵심 타입과 에러를 정의합니다.
// 다른 계층(repository, service, handler)은 모두 이 패키지에 의존합니다.
package domain

import (
	"errors"
	"time"
)

// ErrNotFound는 요청한 Task가 없을 때 쓰는 에러입니다. errors.Is로 비교합니다.
var ErrNotFound = errors.New("할 일을 찾을 수 없습니다")

// ValidationError는 입력값이 규칙에 맞지 않을 때 쓰는 에러입니다. errors.As로 꺼냅니다.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Task는 할 일 하나를 나타냅니다.
type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// UpdateInput은 부분 수정 요청입니다. nil인 필드는 바꾸지 않습니다.
type UpdateInput struct {
	Title *string
	Done  *bool
}
