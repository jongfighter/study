package exercises

import (
	"strings"
	"testing"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		want    string
		wantErr string
	}{
		{"정상", "Go 공부하기", "Go 공부하기", ""},
		{"앞뒤 공백 제거", "  장보기  ", "장보기", ""},
		{"빈 문자열", "", "", "title은 비어 있을 수 없습니다"},
		{"공백만", "   ", "", "title은 비어 있을 수 없습니다"},
		{"한글 100자는 허용", strings.Repeat("가", 100), strings.Repeat("가", 100), ""},
		{"한글 101자는 거부", strings.Repeat("가", 101), "", "title은 100자 이하여야 합니다"},
		{"영문 100자는 허용", strings.Repeat("a", 100), strings.Repeat("a", 100), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateTitle(tt.title)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("예상하지 못한 에러: %v", err)
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
