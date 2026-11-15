package main

import (
	"fmt"
	"sync"
)

func main() {
	counts := make(map[int]int)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range 1000 {
				counts[j] = i
			}
		}()
	}
	wg.Wait()
	fmt.Println("항목 수:", len(counts))
}
