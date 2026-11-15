package repository

import (
	"context"
	"reflect"
	"testing"

	"golab/lesson11/internal/domain"
)

func seedForList(t *testing.T) *Memory {
	t.Helper()
	ctx := context.Background()
	m := NewMemory()
	titles := []string{"Go 공부", "장보기", "go 테스트 작성", "운동", "Go 서버 만들기"}
	for _, title := range titles {
		if _, err := m.Create(ctx, title); err != nil {
			t.Fatal(err)
		}
	}
	done := true
	for _, id := range []int{1, 3} {
		if _, err := m.Update(ctx, id, domain.UpdateInput{Done: &done}); err != nil {
			t.Fatal(err)
		}
	}
	return m
}

func ids(tasks []domain.Task) []int {
	out := make([]int, 0, len(tasks))
	for _, task := range tasks {
		out = append(out, task.ID)
	}
	return out
}

func TestMemoryListFilterAndPagination(t *testing.T) {
	yes, no := true, false
	tests := []struct {
		name      string
		filter    domain.ListFilter
		wantIDs   []int
		wantTotal int
	}{
		{"조건 없음", domain.ListFilter{}, []int{1, 2, 3, 4, 5}, 5},
		{"완료된 것만", domain.ListFilter{Done: &yes}, []int{1, 3}, 2},
		{"완료되지 않은 것만", domain.ListFilter{Done: &no}, []int{2, 4, 5}, 3},
		{"제목 검색은 대소문자를 무시", domain.ListFilter{Query: "go"}, []int{1, 3, 5}, 3},
		{"검색과 완료 여부를 함께", domain.ListFilter{Query: "go", Done: &no}, []int{5}, 1},
		{"limit과 offset", domain.ListFilter{Limit: 2, Offset: 1}, []int{2, 3}, 5},
		{"offset이 끝을 넘으면 빈 목록", domain.ListFilter{Offset: 10}, []int{}, 5},
		{"limit이 남은 개수보다 크면 남은 만큼", domain.ListFilter{Limit: 10, Offset: 3}, []int{4, 5}, 5},
		{"일치하는 것이 없으면 빈 목록", domain.ListFilter{Query: "없는 제목"}, []int{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, err := seedForList(t).List(context.Background(), tt.filter)
			if err != nil {
				t.Fatal(err)
			}
			if got := ids(page.Tasks); !reflect.DeepEqual(got, tt.wantIDs) {
				t.Errorf("ids = %v, want %v", got, tt.wantIDs)
			}
			if page.Total != tt.wantTotal {
				t.Errorf("Total = %d, want %d", page.Total, tt.wantTotal)
			}
			if page.Tasks == nil {
				t.Error("Tasks는 nil이 아니라 빈 슬라이스여야 JSON에서 []로 나옵니다")
			}
		})
	}
}
