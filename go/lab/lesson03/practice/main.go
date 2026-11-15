package main

import (
	"errors"
	"fmt"
	"strconv"
)

func average(nums ...float64) (float64, error) {
	if len(nums) == 0 {
		return 0, errors.New("평균을 구할 값이 없습니다")
	}
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total / float64(len(nums)), nil
}

func fizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return strconv.Itoa(n)
	}
}

func main() {
	avg, err := average(1, 2, 3, 4)
	fmt.Println("평균:", avg, err)

	_, err = average()
	fmt.Println("오류:", err)

	for i := 1; i <= 15; i++ {
		fmt.Println(fizzBuzz(i))
	}
}
