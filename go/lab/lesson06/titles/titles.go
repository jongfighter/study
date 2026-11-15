// Package titles 는 Task 제목을 다듬고 검증합니다.
package titles

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const MaxLength = 50

var (
	ErrEmpty   = errors.New("제목이 비어 있습니다")
	ErrTooLong = errors.New("제목이 너무 깁니다")
)

// Normalize 는 앞뒤 공백을 지우고 연속된 공백을 하나로 줄인 제목을 돌려줍니다.
func Normalize(raw string) (string, error) {
	title := strings.Join(strings.Fields(raw), " ")
	if title == "" {
		return "", ErrEmpty
	}
	if utf8.RuneCountInString(title) > MaxLength {
		return "", ErrTooLong
	}
	return title, nil
}
