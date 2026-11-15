package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	goroutines = 100
	perWorker  = 10000
)

func count(inc func()) {
	var wg sync.WaitGroup
	for range goroutines {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range perWorker {
				inc()
			}
		}()
	}
	wg.Wait()
}

func main() {
	want := goroutines * perWorker

	unsafe := 0
	count(func() { unsafe++ })
	fmt.Printf("잠금 없음: %d (기대값 %d)\n", unsafe, want)

	var mu sync.Mutex
	locked := 0
	count(func() {
		mu.Lock()
		locked++
		mu.Unlock()
	})
	fmt.Printf("Mutex: %d\n", locked)

	var total atomic.Int64
	count(func() { total.Add(1) })
	fmt.Printf("atomic: %d\n", total.Load())
}
