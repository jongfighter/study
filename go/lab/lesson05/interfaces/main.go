package main

import "fmt"

type Notifier interface {
	Notify(message string) string
}

type EmailNotifier struct {
	Address string
}

func (e EmailNotifier) Notify(message string) string {
	return "이메일(" + e.Address + "): " + message
}

type ConsoleNotifier struct{}

func (ConsoleNotifier) Notify(message string) string {
	return "콘솔: " + message
}

type Logger struct {
	Prefix string
}

func (l Logger) Log(message string) {
	fmt.Println(l.Prefix + message)
}

type Service struct {
	Logger
	Notifier
}

type Task struct {
	ID    int
	Title string
}

func (t Task) String() string {
	return fmt.Sprintf("Task#%d(%s)", t.ID, t.Title)
}

func broadcast(notifiers []Notifier, message string) {
	for _, n := range notifiers {
		fmt.Println(n.Notify(message))
	}
}

func describe(value any) string {
	switch v := value.(type) {
	case EmailNotifier:
		return "이메일 알림, 주소=" + v.Address
	case ConsoleNotifier:
		return "콘솔 알림"
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("알 수 없는 타입 %T", v)
	}
}

func main() {
	notifiers := []Notifier{
		EmailNotifier{Address: "dev@example.com"},
		ConsoleNotifier{},
	}
	broadcast(notifiers, "배포가 끝났습니다")

	service := Service{
		Logger:   Logger{Prefix: "[service] "},
		Notifier: ConsoleNotifier{},
	}
	service.Log("시작")
	fmt.Println(service.Notify("임베딩으로 메서드를 물려받습니다"))

	fmt.Println(Task{ID: 1, Title: "인터페이스 익히기"})

	fmt.Println(describe(EmailNotifier{Address: "a@b.c"}))
	fmt.Println(describe(ConsoleNotifier{}))
	fmt.Println(describe(42))
}
