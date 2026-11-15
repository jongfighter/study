package main

import "fmt"

const appName = "golab"

const (
	statusTodo = iota
	statusDoing
	statusDone
)

func main() {
	var count int = 3
	var title string
	name := "Go"
	ratio := 1.5
	done := false

	fmt.Println(count, title == "", name, ratio, done)
	fmt.Printf("타입: %T %T %T %T\n", count, name, ratio, done)

	var (
		zeroInt    int
		zeroFloat  float64
		zeroBool   bool
		zeroString string
		zeroPtr    *int
	)
	fmt.Printf("zero value: %d %v %t %q %v\n", zeroInt, zeroFloat, zeroBool, zeroString, zeroPtr)

	fmt.Println("상수:", appName, statusTodo, statusDoing, statusDone)

	scaled := float64(count) * ratio
	fmt.Println("변환:", scaled)

	word := "한글"
	fmt.Println("바이트 수:", len(word), "글자 수:", len([]rune(word)))
}
