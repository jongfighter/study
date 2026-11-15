// Package exercises는 레슨 08 연습 문제의 정답 예시입니다.
package exercises

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const maxTitleRunes = 100

// ValidateTitle은 제목의 앞뒤 공백을 지우고, 비어 있거나 100글자(룬 기준)를 넘으면 에러를 돌려줍니다.
// len(title)은 글자 수가 아니라 바이트 수이므로 utf8.RuneCountInString으로 셉니다.
func ValidateTitle(title string) (string, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return "", errors.New("title은 비어 있을 수 없습니다")
	}
	if utf8.RuneCountInString(title) > maxTitleRunes {
		return "", errors.New("title은 100자 이하여야 합니다")
	}
	return title, nil
}
