package handler

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"golab/lesson11/internal/domain"
)

func TestParseListFilter(t *testing.T) {
	yes, no := true, false
	tests := []struct {
		name    string
		query   string
		want    domain.ListFilter
		wantErr bool
	}{
		{"아무것도 없음", "", domain.ListFilter{}, false},
		{"done=true", "done=true", domain.ListFilter{Done: &yes}, false},
		{"done=false", "done=false", domain.ListFilter{Done: &no}, false},
		{"모두 지정", "done=true&q=go&limit=10&offset=20", domain.ListFilter{Done: &yes, Query: "go", Limit: 10, Offset: 20}, false},
		{"done이 true/false가 아님", "done=maybe", domain.ListFilter{}, true},
		{"limit이 숫자가 아님", "limit=abc", domain.ListFilter{}, true},
		{"offset이 숫자가 아님", "offset=1.5", domain.ListFilter{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatal(err)
			}
			got, err := parseListFilter(values)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestListTasks(t *testing.T) {
	var gotFilter domain.ListFilter
	list := func(_ context.Context, f domain.ListFilter) (domain.Page, error) {
		gotFilter = f
		return domain.Page{Tasks: []domain.Task{{ID: 3, Title: "Go 공부"}}, Total: 41}, nil
	}
	h := New(fakeTasks{list: list}, testLogger).Routes()

	rec := do(t, h, http.MethodGet, "/tasks?q=go&limit=1&offset=2", "")

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d (body: %s)", rec.Code, rec.Body)
	}
	if got := rec.Header().Get("X-Total-Count"); got != "41" {
		t.Errorf("X-Total-Count = %q, want 41", got)
	}
	if !strings.Contains(rec.Body.String(), `"id":3`) {
		t.Errorf("body = %s", rec.Body)
	}
	if gotFilter.Query != "go" || gotFilter.Limit != 1 || gotFilter.Offset != 2 {
		t.Errorf("서비스로 넘어간 조건이 다릅니다: %+v", gotFilter)
	}
}

func TestListTasksRejectsBadQuery(t *testing.T) {
	list := func(context.Context, domain.ListFilter) (domain.Page, error) {
		t.Error("잘못된 쿼리면 서비스를 부르면 안 됩니다")
		return domain.Page{}, nil
	}
	h := New(fakeTasks{list: list}, testLogger).Routes()

	rec := do(t, h, http.MethodGet, "/tasks?done=maybe", "")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}

func TestListTasksMapsValidationError(t *testing.T) {
	list := func(context.Context, domain.ListFilter) (domain.Page, error) {
		return domain.Page{}, &domain.ValidationError{Message: "limit은 1 이상 100 이하여야 합니다"}
	}
	h := New(fakeTasks{list: list}, testLogger).Routes()

	rec := do(t, h, http.MethodGet, "/tasks?limit=1000", "")

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", rec.Code)
	}
}
