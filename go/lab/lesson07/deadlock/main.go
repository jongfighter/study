package main

import "fmt"

func main() {
	ch := make(chan int)
	ch <- 1
	fmt.Println("여기에는 도달하지 못합니다")
}
