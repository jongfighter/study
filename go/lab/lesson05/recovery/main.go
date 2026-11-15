package main

import "fmt"

func deferOrder() {
	for i := 1; i <= 3; i++ {
		defer fmt.Println("defer", i)
	}
	fmt.Println("본문 끝")
}

func safeDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("복구했습니다: %v", r)
		}
	}()
	return a / b, nil
}

func main() {
	deferOrder()

	result, err := safeDivide(10, 2)
	fmt.Println(result, err)

	result, err = safeDivide(1, 0)
	fmt.Println(result, err)

	fmt.Println("프로그램은 계속 실행됩니다")
}
