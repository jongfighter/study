package main

import "fmt"

type Task struct {
	ID    int
	Title string
	Done  bool
}

func NewTask(id int, title string) *Task {
	return &Task{ID: id, Title: title}
}

func (t Task) Describe() string {
	mark := " "
	if t.Done {
		mark = "x"
	}
	return fmt.Sprintf("[%s] #%d %s", mark, t.ID, t.Title)
}

func (t *Task) Toggle() {
	t.Done = !t.Done
}

func (t Task) BrokenToggle() {
	t.Done = !t.Done
}

func main() {
	task := Task{ID: 1, Title: "Go 공부하기"}
	fmt.Println(task.Describe())
	fmt.Printf("%v\n%+v\n", task, task)

	task.BrokenToggle()
	fmt.Println("BrokenToggle 이후:", task.Done)

	task.Toggle()
	fmt.Println("Toggle 이후:", task.Done)

	copied := task
	copied.Title = "복사본"
	fmt.Println(task.Title, "/", copied.Title)

	pointer := &task
	pointer.Title = "포인터로 수정"
	fmt.Println(task.Title)

	created := NewTask(2, "테스트 작성")
	created.Toggle()
	fmt.Println(created.Describe())

	tasks := []Task{{ID: 3, Title: "배포하기"}, {ID: 4, Title: "회고하기"}}
	for _, t := range tasks {
		t.Done = true
	}
	fmt.Println("range 복사본을 바꾼 뒤:", tasks[0].Done, tasks[1].Done)

	for i := range tasks {
		tasks[i].Done = true
	}
	fmt.Println("인덱스로 바꾼 뒤:", tasks[0].Done, tasks[1].Done)

	var zero Task
	fmt.Printf("zero value: %+v\n", zero)
	fmt.Println(task == Task{ID: 1, Title: "포인터로 수정", Done: true})
}
