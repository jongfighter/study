package main

import (
	"errors"
	"fmt"
)

func add(a, b int) int {
	return a + b
}

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("0으로 나눌 수 없습니다")
	}
	return a / b, nil
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func minMax(nums ...int) (lowest, highest int) {
	lowest, highest = nums[0], nums[0]
	for _, n := range nums {
		if n < lowest {
			lowest = n
		}
		if n > highest {
			highest = n
		}
	}
	return
}

func makeCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	fmt.Println("add(2, 3) =", add(2, 3))

	quotient, err := divide(7, 2)
	fmt.Println("7 / 2 =", quotient, err)

	_, err = divide(1, 0)
	fmt.Println("오류:", err)

	values := []int{4, 5, 6}
	fmt.Println("합계:", sum(), sum(1, 2, 3), sum(values...))

	lowest, highest := minMax(3, 9, 1)
	fmt.Println("최소/최대:", lowest, highest)

	next := makeCounter()
	fmt.Println("카운터:", next(), next(), next())
}
