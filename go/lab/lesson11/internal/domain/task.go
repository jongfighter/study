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

// ListFilter는 목록 조회 조건입니다.
type ListFilter struct {
	Done   *bool  // nil이면 완료 여부와 상관없이 모두 조회합니다.
	Query  string // 비어 있지 않으면 제목에 이 문자열이 들어 있는(대소문자 무시) Task만 조회합니다.
	Limit  int    // 한 번에 돌려줄 최대 개수입니다. 0이면 서비스의 기본값을 씁니다.
	Offset int    // 앞에서부터 건너뛸 개수입니다.
}

// Page는 목록 조회 결과 한 쪽입니다.
type Page struct {
	Tasks []Task
	Total int // 조건에 맞는 전체 개수(Limit, Offset을 적용하기 전)
}
