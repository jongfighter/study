package main

import "fmt"

func main() {
	score := 82
	if score >= 90 {
		fmt.Println("등급 A")
	} else if score >= 80 {
		fmt.Println("등급 B")
	} else {
		fmt.Println("등급 C 이하")
	}

	if n := len("Go"); n > 1 {
		fmt.Println("if 안에서만 보이는 n =", n)
	}

	level := "바깥"
	if score > 0 {
		level := "안쪽"
		fmt.Println("블록 안:", level)
	}
	fmt.Println("블록 밖:", level)

	sum := 0
	for i := 1; i <= 5; i++ {
		sum += i
	}
	fmt.Println("1부터 5까지 합:", sum)

	n := 1
	for n < 100 {
		n *= 2
	}
	fmt.Println("100을 넘는 첫 2의 거듭제곱:", n)

	fmt.Print("range 3:")
	for i := range 3 {
		fmt.Print(" ", i)
	}
	fmt.Println()

	day := "토"
	switch day {
	case "토", "일":
		fmt.Println("주말")
	default:
		fmt.Println("평일")
	}

	switch {
	case score >= 90:
		fmt.Println("우수")
	case score >= 70:
		fmt.Println("보통")
	default:
		fmt.Println("노력 필요")
	}
}
