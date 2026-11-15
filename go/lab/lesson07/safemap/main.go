package main

import (
	"fmt"
	"sync"
)

type safeCounts struct {
	mu     sync.Mutex
	counts map[int]int
}

func (s *safeCounts) set(key, value int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key] = value
}

func main() {
	s := &safeCounts{counts: make(map[int]int)}

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 1000 {
				s.set(j, i)
			}
		}()
	}
	wg.Wait()
	fmt.Println("항목 수:", len(s.counts))
}
