// Package service는 Task에 관한 업무 규칙을 담습니다. HTTP도, 저장 방식도 모릅니다.
package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"golab/lesson11/internal/domain"
)

const (
	maxTitleRunes = 100
	defaultLimit  = 20
	maxLimit      = 100
)

// Repository는 서비스가 필요로 하는 저장소 동작입니다.
// 인터페이스는 구현하는 쪽이 아니라 사용하는 쪽(service)에서 정의합니다.
type Repository interface {
	List(ctx context.Context, f domain.ListFilter) (domain.Page, error)
	Create(ctx context.Context, title string) (domain.Task, error)
	Get(ctx context.Context, id int) (domain.Task, error)
	Update(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error)
	Delete(ctx context.Context, id int) error
}

// Service는 Task 업무 규칙을 실행합니다.
type Service struct {
	repo Repository
}

// New는 repo를 주입받아 Service를 만듭니다.
func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// List는 조건을 검증하고 기본값을 채운 뒤 저장소에서 한 쪽을 가져옵니다.
func (s *Service) List(ctx context.Context, f domain.ListFilter) (domain.Page, error) {
	if f.Limit == 0 {
		f.Limit = defaultLimit
	}
	if f.Limit < 1 || f.Limit > maxLimit {
		return domain.Page{}, &domain.ValidationError{Message: "limit은 1 이상 100 이하여야 합니다"}
	}
	if f.Offset < 0 {
		return domain.Page{}, &domain.ValidationError{Message: "offset은 0 이상이어야 합니다"}
	}
	f.Query = strings.TrimSpace(f.Query)
	return s.repo.List(ctx, f)
}

func (s *Service) Create(ctx context.Context, title string) (domain.Task, error) {
	title, err := normalizeTitle(title)
	if err != nil {
		return domain.Task{}, err
	}
	return s.repo.Create(ctx, title)
}

func (s *Service) Get(ctx context.Context, id int) (domain.Task, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int, in domain.UpdateInput) (domain.Task, error) {
	if in.Title == nil && in.Done == nil {
		return domain.Task{}, &domain.ValidationError{Message: "수정할 필드(title, done)가 없습니다"}
	}
	if in.Title != nil {
		title, err := normalizeTitle(*in.Title)
		if err != nil {
			return domain.Task{}, err
		}
		in.Title = &title
	}
	return s.repo.Update(ctx, id, in)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func normalizeTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", &domain.ValidationError{Message: "title은 비어 있을 수 없습니다"}
	}
	if utf8.RuneCountInString(title) > maxTitleRunes {
		return "", &domain.ValidationError{Message: "title은 100자 이하여야 합니다"}
	}
	return title, nil
}
