package exercises

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAPIKey(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	protected := RequireAPIKey("secret")(ok)

	tests := []struct {
		name       string
		method     string
		path       string
		header     string
		wantStatus int
	}{
		{"키가 맞으면 통과", http.MethodGet, "/tasks", "secret", http.StatusOK},
		{"키가 없으면 401", http.MethodGet, "/tasks", "", http.StatusUnauthorized},
		{"키가 틀리면 401", http.MethodPost, "/tasks", "wrong", http.StatusUnauthorized},
		{"health는 키 없이 통과", http.MethodGet, "/health", "", http.StatusOK},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			if tt.header != "" {
				req.Header.Set("X-API-Key", tt.header)
			}
			rec := httptest.NewRecorder()

			protected.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
