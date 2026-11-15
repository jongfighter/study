package workerpool

import (
	"context"
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"
)

func double(_ context.Context, n int) (int, error) {
	return n * 2, nil
}

func TestMapResults(t *testing.T) {
	tests := []struct {
		name    string
		workers int
		nums    []int
		want    []int
	}{
		{"빈 입력", 3, []int{}, []int{}},
		{"입력 순서 유지", 3, []int{1, 2, 3, 4, 5, 6, 7, 8}, []int{2, 4, 6, 8, 10, 12, 14, 16}},
		{"작업보다 워커가 많음", 10, []int{5, 6}, []int{10, 12}},
		{"워커 0개는 1개로 취급", 0, []int{1, 2, 3}, []int{2, 4, 6}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Map(context.Background(), tt.workers, tt.nums, double)
			if err != nil {
				t.Fatalf("예상하지 못한 에러: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapFirstErrorCancelsOthers(t *testing.T) {
	boom := errors.New("boom")
	fn := func(ctx context.Context, n int) (int, error) {
		if n == 3 {
			return 0, boom
		}
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(2 * time.Second):
			return n, nil
		}
	}

	start := time.Now()
	_, err := Map(context.Background(), 4, []int{1, 2, 3, 4}, fn)
	if !errors.Is(err, boom) {
		t.Fatalf("got %v, want boom", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("나머지 작업이 취소되지 않았습니다: %v 걸림", elapsed)
	}
}

func TestMapParentCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := Map(ctx, 2, []int{1, 2, 3}, double)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("got %v, want context.Canceled", err)
	}
}

func TestMapLimitsWorkers(t *testing.T) {
	const workers = 3
	var active, peak atomic.Int32
	fn := func(_ context.Context, n int) (int, error) {
		now := active.Add(1)
		for {
			old := peak.Load()
			if now <= old || peak.CompareAndSwap(old, now) {
				break
			}
		}
		time.Sleep(10 * time.Millisecond)
		active.Add(-1)
		return n, nil
	}

	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if _, err := Map(context.Background(), workers, nums, fn); err != nil {
		t.Fatal(err)
	}
	if got := peak.Load(); got > workers {
		t.Errorf("동시에 실행된 작업 수 %d가 워커 수 %d보다 많습니다", got, workers)
	}
}
