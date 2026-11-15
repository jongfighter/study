package main

import "fmt"

func generate(n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- i
		}
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			out <- v * v
		}
	}()
	return out
}

func main() {
	for v := range square(generate(5)) {
		fmt.Println("제곱:", v)
	}

	buffered := make(chan string, 2)
	buffered <- "첫째"
	buffered <- "둘째"
	fmt.Println("버퍼 채널 len/cap:", len(buffered), cap(buffered))
	fmt.Println(<-buffered, <-buffered)
}
