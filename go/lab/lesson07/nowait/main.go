package main

import "fmt"

func main() {
	go fmt.Println("고루틴에서 출력")
	fmt.Println("main 종료")
}
