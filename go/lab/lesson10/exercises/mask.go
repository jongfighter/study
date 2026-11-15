package exercises

import (
	"log/slog"
	"strings"
)

var secretKeys = map[string]bool{
	"password":      true,
	"token":         true,
	"authorization": true,
}

// MaskSecrets는 slog.HandlerOptions.ReplaceAttr에 넘겨 쓰는 함수입니다.
// 키가 password, token, authorization(대소문자 무시)인 속성의 값을 ***로 바꿉니다.
func MaskSecrets(groups []string, a slog.Attr) slog.Attr {
	if secretKeys[strings.ToLower(a.Key)] {
		return slog.String(a.Key, "***")
	}
	return a
}
