package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func producer(ctx context.Context, out chan<- int) error {
	for i := 1; ; i++ {
		select {
		case out <- i:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func cancelDemo() {
	ctx, cancel := context.WithCancel(context.Background())
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := producer(ctx, out)
		fmt.Println("생산자 종료:", err)
	}()

	for range 3 {
		fmt.Println("받음:", <-out)
	}
	fmt.Println("cancel 호출")
	cancel()
	wg.Wait()
}

func timeoutDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("작업 완료")
	case <-ctx.Done():
		fmt.Println("중단:", ctx.Err())
	}
}

func main() {
	cancelDemo()
	timeoutDemo()
}
