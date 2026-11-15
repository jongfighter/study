package titles

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestNormalizeTrimsSpaces(t *testing.T) {
	got, err := Normalize("  장보기  ")
	if err != nil {
		t.Fatalf("예상하지 못한 오류: %v", err)
	}
	if got != "장보기" {
		t.Errorf("Normalize() = %q; %q 이어야 합니다", got, "장보기")
	}
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{"앞뒤 공백", "  장보기  ", "장보기", nil},
		{"연속 공백", "Go   공부하기", "Go 공부하기", nil},
		{"탭과 줄바꿈", "a\tb\nc", "a b c", nil},
		{"빈 문자열", "", "", ErrEmpty},
		{"공백만", "   ", "", ErrEmpty},
		{"최대 길이", strings.Repeat("가", MaxLength), strings.Repeat("가", MaxLength), nil},
		{"너무 김", strings.Repeat("가", MaxLength+1), "", ErrTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Normalize(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("오류 = %v; %v 이어야 합니다", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q; %q 이어야 합니다", tt.input, got, tt.want)
			}
		})
	}
}

func ExampleNormalize() {
	title, _ := Normalize("  Go   공부하기  ")
	fmt.Println(title)
	// Output: Go 공부하기
}

func BenchmarkNormalize(b *testing.B) {
	for b.Loop() {
		Normalize("  Go   공부하기  ")
	}
}
