package main

import (
	"maps"
	"slices"
	"testing"
)

func TestWordCount(t *testing.T) {
	got := wordCount("I ate a donut. Then I ate another donut.")
	want := map[string]int{"I": 2, "ate": 2, "a": 1, "donut.": 2, "Then": 1, "another": 1}
	if !maps.Equal(got, want) {
		t.Errorf("wordCount = %v; %v 이어야 합니다", got, want)
	}
}

func TestPendingTitles(t *testing.T) {
	tasks := []Task{{ID: 1, Title: "a", Done: true}, {ID: 2, Title: "b"}, {ID: 3, Title: "c"}}
	if got := pendingTitles(tasks); !slices.Equal(got, []string{"b", "c"}) {
		t.Errorf("pendingTitles = %v", got)
	}
}

func TestCompleteAll(t *testing.T) {
	tasks := []Task{{ID: 1, Title: "a"}, {ID: 2, Title: "b"}}
	completeAll(tasks)
	if got := pendingTitles(tasks); len(got) != 0 {
		t.Errorf("모두 끝내야 하는데 남은 일: %v", got)
	}
	if !tasks[0].Done || !tasks[1].Done {
		t.Error("원본 슬라이스의 요소가 바뀌어야 합니다")
	}
}
