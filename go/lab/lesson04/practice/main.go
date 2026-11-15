package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

func wordCount(text string) map[string]int {
	counts := make(map[string]int)
	for _, word := range strings.Fields(text) {
		counts[word]++
	}
	return counts
}

func pendingTitles(tasks []Task) []string {
	var titles []string
	for _, task := range tasks {
		if !task.Done {
			titles = append(titles, task.Title)
		}
	}
	return titles
}

func completeAll(tasks []Task) {
	for i := range tasks {
		tasks[i].Done = true
	}
}

func main() {
	counts := wordCount("I ate a donut. Then I ate another donut.")
	for _, word := range slices.Sorted(maps.Keys(counts)) {
		fmt.Println(word, counts[word])
	}

	tasks := []Task{
		{ID: 1, Title: "글쓰기", Done: true},
		{ID: 2, Title: "빨래하기"},
		{ID: 3, Title: "운동하기"},
	}
	fmt.Println("남은 일:", pendingTitles(tasks))

	completeAll(tasks)
	fmt.Println("모두 끝낸 뒤:", pendingTitles(tasks), len(pendingTitles(tasks)))
}
