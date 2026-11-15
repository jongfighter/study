package main

import (
	"fmt"
	"time"
)

func slowWork() <-chan string {
	result := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond)
		result <- "작업 결과"
	}()
	return result
}

func waitFor(timeout time.Duration) {
	select {
	case r := <-slowWork():
		fmt.Println("성공:", r)
	case <-time.After(timeout):
		fmt.Println("시간 초과:", timeout)
	}
}

func main() {
	waitFor(50 * time.Millisecond)
	waitFor(500 * time.Millisecond)
}
