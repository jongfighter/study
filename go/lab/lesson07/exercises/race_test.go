package exercises

import (
	"context"
	"errors"
	"testing"
	"time"
)

func after(d time.Duration, answer string, canceled chan<- struct{}) Source {
	return func(ctx context.Context) string {
		select {
		case <-time.After(d):
			return answer
		case <-ctx.Done():
			if canceled != nil {
				canceled <- struct{}{}
			}
			return ""
		}
	}
}

func TestRaceFastestWins(t *testing.T) {
	canceled := make(chan struct{}, 1)
	got, err := Race(context.Background(), time.Second,
		after(300*time.Millisecond, "느림", canceled),
		after(10*time.Millisecond, "빠름", nil),
	)
	if err != nil {
		t.Fatal(err)
	}
	if got != "빠름" {
		t.Errorf("got %q, want %q", got, "빠름")
	}

	select {
	case <-canceled:
	case <-time.After(time.Second):
		t.Error("늦은 작업에 취소 신호가 전달되지 않았습니다")
	}
}

func TestRaceTimeout(t *testing.T) {
	_, err := Race(context.Background(), 30*time.Millisecond,
		after(500*time.Millisecond, "느림1", nil),
		after(500*time.Millisecond, "느림2", nil),
	)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("got %v, want context.DeadlineExceeded", err)
	}
}
