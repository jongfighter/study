package exercises

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func newMaskedLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buf, &slog.HandlerOptions{ReplaceAttr: MaskSecrets}))
}

func TestMaskSecretsHidesSensitiveValues(t *testing.T) {
	var buf bytes.Buffer
	logger := newMaskedLogger(&buf)

	logger.Info("로그인 시도", "user", "gopher", "Password", "p@ssw0rd", "token", "abc123")

	out := buf.String()
	for _, secret := range []string{"p@ssw0rd", "abc123"} {
		if strings.Contains(out, secret) {
			t.Errorf("로그에 비밀 값 %q가 그대로 남았습니다: %s", secret, out)
		}
	}
	for _, want := range []string{"user=gopher", "Password=***", "token=***"} {
		if !strings.Contains(out, want) {
			t.Errorf("로그에 %q가 있어야 합니다: %s", want, out)
		}
	}
}

func TestMaskSecretsKeepsOtherAttributes(t *testing.T) {
	var buf bytes.Buffer
	logger := newMaskedLogger(&buf)

	logger.Warn("요청", "path", "/tasks", "status", 404)

	out := buf.String()
	if !strings.Contains(out, "path=/tasks") || !strings.Contains(out, "status=404") {
		t.Errorf("다른 속성은 그대로여야 합니다: %s", out)
	}
}
