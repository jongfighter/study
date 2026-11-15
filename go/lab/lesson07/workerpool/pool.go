// Package workerpool은 고정된 수의 고루틴으로 숫자 작업을 나눠 처리하는 예제입니다.
package workerpool

import (
	"context"
	"sync"
)

// Func는 숫자 하나를 처리하는 작업입니다. ctx가 취소되면 가능한 한 빨리 돌아와야 합니다.
type Func func(ctx context.Context, n int) (int, error)

// Map은 workers개의 고루틴으로 nums를 처리하고 입력 순서대로 결과를 돌려줍니다.
// 작업 하나가 에러를 내면 나머지를 취소하고 처음 발생한 에러를 돌려줍니다.
func Map(ctx context.Context, workers int, nums []int, fn Func) ([]int, error) {
	if workers < 1 {
		workers = 1
	}
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]int, len(nums))
	jobs := make(chan int)

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				v, err := fn(runCtx, nums[i])
				if err != nil {
					mu.Lock()
					if firstErr == nil {
						firstErr = err
						cancel()
					}
					mu.Unlock()
					continue
				}
				results[i] = v
			}
		}()
	}

feed:
	for i := range nums {
		select {
		case jobs <- i:
		case <-runCtx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
