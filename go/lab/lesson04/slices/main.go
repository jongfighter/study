package main

import "fmt"

func double(values []int) {
	for i := range values {
		values[i] *= 2
	}
}

func main() {
	numbers := []int{1, 2, 3}
	fmt.Println(numbers, len(numbers), cap(numbers))

	numbers = append(numbers, 4)
	fmt.Println(numbers, len(numbers), cap(numbers))

	base := []int{10, 20, 30, 40, 50}
	part := base[1:3]
	fmt.Println("part:", part, len(part), cap(part))

	part[0] = 99
	fmt.Println("base:", base)

	part = append(part, 77)
	fmt.Println("base:", base)

	independent := append([]int(nil), base[1:3]...)
	independent[0] = -1
	fmt.Println("복사본:", independent, "원본:", base)

	var empty []int
	fmt.Println("nil 슬라이스:", empty == nil, len(empty))
	empty = append(empty, 1)
	fmt.Println(empty)

	scores := []int{1, 2, 3}
	double(scores)
	fmt.Println("double 이후:", scores)

	array := [3]int{1, 2, 3}
	arrayCopy := array
	arrayCopy[0] = 100
	fmt.Println("배열:", array, arrayCopy)
}
